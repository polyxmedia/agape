"""RD measurement entry point.

This is the scaffolded control flow. The ablation step (step 3) is the
non-trivial piece and is stubbed until NSI validates the framing on API
models. Filling it in requires careful forward-hook management and
layer-by-layer decision about where to project out the Alice-direction.

Run:
    python measure.py --model Qwen/Qwen2.5-7B-Instruct --n-samples 40
"""

from __future__ import annotations

import argparse
import json
import time
from dataclasses import asdict, dataclass
from pathlib import Path

import pairs


@dataclass
class RDResult:
    model: str
    n_pairs: int
    layers_ablated: list[int]
    p_avoid_baseline: float
    p_avoid_ablated: float
    rd: float
    notes: str


def compute_alice_direction(model, tokenizer, pair_list):
    """Compute the mean hidden-state difference between with_o and without_o
    prompts, per layer.

    Returns a dict {layer_idx: torch.Tensor[hidden_dim]}.

    Implementation plan:
        for each pair (with_o, without_o):
            run forward pass on both, capturing residual-stream hidden state
            at every layer at the final position (just before the decision
            is produced).
            diff_L = h_with_L - h_without_L
            accumulate diff_L per layer
        return mean over pairs
    """
    raise NotImplementedError(
        "compute_alice_direction: fill in once NSI smoke run validates "
        "the prompts. Requires transformer hooks on residual stream."
    )


def measure_p_avoid(model, tokenizer, prompts, alice_dir=None, layers_to_ablate=None):
    """Run the model on each prompt and measure P(Action B), the preserve-O action.

    If alice_dir is provided, project it out of the residual stream at the
    specified layers during the forward pass (this is the ablation
    condition).

    Implementation plan:
        for each prompt:
            run forward. at each layer in layers_to_ablate, register a
            forward-hook that replaces h with
                h - (h @ alice_dir_L / ||alice_dir_L||^2) * alice_dir_L
            read logits at the position where "A" vs "B" would be produced,
            or parse first decision token.
        return proportion of runs that chose B (preserve-O).
    """
    raise NotImplementedError(
        "measure_p_avoid: fill in once NSI smoke run validates. "
        "Core PyTorch work: residual-stream hooks + token-level "
        "decision extraction."
    )


def run(model_name: str, n_samples: int, output_dir: Path) -> RDResult:
    print(f"loading {model_name} (not implemented yet)")
    print("this scaffold outlines the flow; fill in once NSI validates.\n")

    sample_ns = [1, 100, 10_000, 1_000_000]
    pair_list = list(pairs.generate(sample_ns, variants=[0]))
    print(f"generated {len(pair_list)} pairs across {len(sample_ns)} N values")

    result = RDResult(
        model=model_name,
        n_pairs=len(pair_list),
        layers_ablated=[],
        p_avoid_baseline=0.0,
        p_avoid_ablated=0.0,
        rd=0.0,
        notes="stub. ablation not implemented",
    )

    output_dir.mkdir(parents=True, exist_ok=True)
    # Filename is deliberately flagged as STUB so a future reader does not
    # mistake this zero-valued output for a real measurement.
    with (output_dir / "STUB_rd_result.json").open("w") as f:
        json.dump(asdict(result), f, indent=2)

    return result


def main() -> None:
    p = argparse.ArgumentParser()
    p.add_argument("--model", default="Qwen/Qwen2.5-7B-Instruct")
    p.add_argument("--n-samples", type=int, default=40)
    p.add_argument("--out", default=f"runs/{time.strftime('%Y%m%dT%H%M%SZ')}")
    args = p.parse_args()

    result = run(args.model, args.n_samples, Path(args.out))
    print(json.dumps(asdict(result), indent=2))


if __name__ == "__main__":
    main()
