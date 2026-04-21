"""Verify that rd/pairs.py mirrors the Go-side scenarios.

The RD experiment relies on both substrates rendering the same prompts. If
they diverge, the 'with O' Python prompt no longer corresponds to any
prompt a subject model sees on the Go side, and RD is measuring a
different distribution than NSI.

This script checks two things:

1. ALICE_BIO matches byte-for-byte across Go and Python.
2. The three forced_choice templates, after normalising whitespace, match
   across Go and Python (parameterised on NAME and BIO on both sides).

Exit 0 on match; exit 1 with a diff on mismatch.
"""

from __future__ import annotations

import difflib
import re
import sys
from pathlib import Path

import pairs


GO_FILE = Path(__file__).resolve().parent.parent / "internal" / "scenarios" / "scenarios.go"


def _extract_const_string(name: str, text: str) -> str:
    """Pull a backtick-delimited Go const (possibly multi-segment) into one string."""
    match = re.search(rf"const\s+{name}\s*=\s*((?:`[^`]*`\s*\+?\s*)+)", text)
    if not match:
        raise RuntimeError(
            f"could not locate {name} const in {GO_FILE}. "
            "the Go declaration shape may have changed; update check_drift.py"
        )
    segments = re.findall(r"`([^`]*)`", match.group(1))
    return "".join(segments)


def _extract_forced_choice_templates(text: str) -> list[str]:
    """Pull the three forced_choice templates out of the Go templates map.

    Each template is a backtick-delimited string (possibly concatenated
    across multiple lines with `+`). Returns them in source order.
    """
    match = re.search(
        r'"forced_choice":\s*\{(.*?)\n\t\},\n\}',
        text,
        re.DOTALL,
    )
    if not match:
        raise RuntimeError(
            "could not locate forced_choice templates block in "
            f"{GO_FILE}. update check_drift.py"
        )
    block = match.group(1)
    # Each template begins with a `, ends with `, with possible + between
    # segments. Greedy across backticks would merge templates, so walk
    # segment-by-segment and split when we see a `,` or `,\n` at the top level.
    templates: list[str] = []
    current: list[str] = []
    depth = 0
    i = 0
    while i < len(block):
        c = block[i]
        if c == "`":
            end = block.index("`", i + 1)
            current.append(block[i + 1 : end])
            i = end + 1
            continue
        if c == "," and depth == 0 and current:
            templates.append("".join(current))
            current = []
        i += 1
    if current:
        templates.append("".join(current))
    return templates


def _normalise(s: str) -> str:
    return re.sub(r"\s+", " ", s).strip()


def main() -> int:
    text = GO_FILE.read_text()
    ok = True

    # rd/pairs.py's ALICE_BIO is the rich bio. Go has two bios now
    # (AliceBioRich, AliceBioNeutral); Python only mirrors the rich one.
    go_bio = _extract_const_string("AliceBioRich", text)
    if go_bio != pairs.ALICE_BIO:
        ok = False
        print("DRIFT: ALICE_BIO differs between Go and Python", file=sys.stderr)
        sys.stderr.writelines(
            difflib.unified_diff(
                go_bio.splitlines(keepends=True),
                pairs.ALICE_BIO.splitlines(keepends=True),
                fromfile="go:AliceBio",
                tofile="py:ALICE_BIO",
            )
        )

    go_tpls = _extract_forced_choice_templates(text)
    if len(go_tpls) != len(pairs.FORCED_CHOICE):
        ok = False
        print(
            f"DRIFT: forced_choice template count differs "
            f"(go={len(go_tpls)}, py={len(pairs.FORCED_CHOICE)})",
            file=sys.stderr,
        )
    else:
        for i, (g, p) in enumerate(zip(go_tpls, pairs.FORCED_CHOICE)):
            if _normalise(g) != _normalise(p):
                ok = False
                print(f"DRIFT: forced_choice variant {i} differs", file=sys.stderr)
                sys.stderr.writelines(
                    difflib.unified_diff(
                        _normalise(g).splitlines(keepends=True),
                        _normalise(p).splitlines(keepends=True),
                        fromfile=f"go:variant[{i}]",
                        tofile=f"py:FORCED_CHOICE[{i}]",
                    )
                )

    if ok:
        print("OK: Go and Python agree on ALICE_BIO and forced_choice templates")
        return 0
    return 1


if __name__ == "__main__":
    sys.exit(main())
