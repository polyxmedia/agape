# Contributing to Agape

Thanks for taking an interest. This project is the operationalisation of
the Agape Protocol described in `all-you-need-is-love.md`. Contributions
that move the measurement infrastructure forward, sharpen the
methodology, or extend the empirical coverage are welcome. Contributions
that add scope without adding measurement value will be politely
declined.

## Before you start

Read these three files in order:

1. `all-you-need-is-love.md`. The paper. Explains what the project is
   trying to measure and why behavioural alignment is insufficient.
2. `ANALYSIS.md`. The pre-registered analysis plan. Defines the primary
   statistic, validity gates, and interpretation ladder. New work that
   reports NSI numbers must conform to this plan or explicitly amend it
   in a documented diff.
3. `RESULTS.md`. The first cross-model NSI measurement and the first RD
   measurement. Explains what is currently established and what is not.

## What good contributions look like

**Yes, please:**

- New measurement axes that close gaps the paper identifies (especially
  SpI, the specificity index).
- New subjects and bios for the existing NSI sweep, with the
  pre-registered validity gates respected.
- Layer-by-layer attribution for the RD axis, separating where the
  representational dependence originates.
- Cleaner counterfactuals for the RD pair generator (current Sarah Kim
  pair matches on biographical concreteness; better controls welcome).
- Bug fixes with regression tests.
- Performance improvements that do not change the measurement.
- Documentation that helps the next person understand a non-obvious
  design choice.
- Replication of the published results on different hardware or with
  re-implementations in other languages.

**Please don't:**

- Add features that the paper's argument doesn't motivate. Scope creep
  is the enemy here.
- Train models against the metrics. The metrics are intentionally
  holdout diagnostics. Goodhart applies sharpest in this domain.
- Open issues asking us to report only one of the six break-point
  values. They are reported together for honesty about uncertainty.
- Submit results from sweeps that fail the validity gates without
  flagging them as exploratory per ANALYSIS.md.
- Rewrite the paper or the analysis plan in PRs that also touch code.
  Methodology changes go in their own PR with a written rationale.

## Workflow

1. Open an issue describing what you want to do before writing code,
   especially if the change is non-trivial. This avoids you doing work
   we'd ask you to revert.
2. Fork the repo. Branch from `main`.
3. Code change comes with tests if it touches the Go side
   (`internal/nsi/*`, `internal/scenarios/*`). RD-side changes come with
   at least one verification run committed under `runs/`.
4. Run `go fmt ./... && go vet ./... && go test ./...` before pushing.
5. Run `python3 rd/check_drift.py` before pushing if you touched
   prompts on either substrate.
6. Open a PR. Describe what you measured, what changed, and which
   validity gates the change is meant to address or affect.

## Code style

- Go: standard `gofmt` plus `go vet` clean. Public types documented
  briefly. No comments restating what the code does. Comments only when
  the *why* is non-obvious.
- Python: `ruff` and `mypy` clean if you have them; otherwise just
  Python 3.10+ type hints and reasonable PEP 8.
- No em-dashes anywhere in code, comments, or documentation.
- No "this PR" or "added for X" comments. PR descriptions belong in PR
  descriptions.

## Reporting empirical results

If a PR adds or modifies a measurement, it must include:

- The exact config used (`config.resolved.yaml` from the run).
- The trial-level data (`trials.jsonl`).
- A brief writeup in PR description: subject, bio, N range, primary
  statistic, validity-gate status. If any gate failed, say so.
- A statement of what the result does and does not establish, in the
  paper's terms.

## License

By submitting a contribution you agree it is licensed under the Apache
License, Version 2.0 (see LICENSE), under the same terms as the rest of
the codebase. The paper text is separately licensed CC BY 4.0; PRs that
modify the paper text follow that license.

## Conduct

Argue with the work. Be civil with the people. The methodology is the
target; everything else is process.
