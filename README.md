# agape

Operationalisation of the Agape Protocol described in
[`all-you-need-is-love.md`](all-you-need-is-love.md). The paper argues
that current alignment is insufficient because it shapes behaviour
rather than the structural property that determines behaviour under
pressure. This codebase builds the *measurement* of that property, on
the principle that you cannot iterate toward something you cannot
detect.

## Status

Two of three measurement axes are built and have produced first
results.

| Axis | What it measures | Status |
|---|---|---|
| **NSI** (Non-Substitutability Index) | Behavioural break-point under increasing instrumental gain. Operationalises Condition 1 of the paper. | Shipped. First cross-model measurement: `RESULTS.md`. |
| **RD** (Representational Dependence) | Whether harm-avoidance is causally tied to the model's representation of the named O. Operationalises Condition 2 at the representational level. | MVP shipped. First measurement on Qwen 2.5 7B: `RESULTS.md`. |
| **SpI** (Specificity Index) | Whether the avoidance is for the specified individual or for named individuals in general. | Planned. See `ROADMAP.md`. |

## Headline result so far

Two production frontier models differ by approximately four orders of
magnitude on NSI under identical framing. claude-haiku-4-5's
sustained-confident break-point is not reached within the primary range
(NSI >= 5). gpt-4o-mini's is at N=5 (NSI = 0.699). On the RD side, Qwen
2.5 7B Instruct shows representational dependence consistent with the
Condition 2 prediction at the regime where it actually preserves
(corrected RD = +0.80 at N=1, against a random-direction control of
−0.13).

Full writeup including validity gates, exclusions, and what these
results do and do not establish: [`RESULTS.md`](RESULTS.md).

## Run

```sh
git clone https://github.com/polyxmedia/agape.git
cd agape
cp .env.example .env       # edit with your keys
./run.sh -smoke            # validate the pipeline (~12 API calls)
./run.sh                   # full sweep (~3,120 API calls, 40-60 min)
```

For RD:

```sh
cd rd
python3 -m venv .venv
.venv/bin/pip install -r requirements.txt
.venv/bin/python check_drift.py     # verify Go/Python prompt drift
.venv/bin/python measure.py         # ~5 min on M-series Mac with MPS
```

## Layout

```
all-you-need-is-love.md     The paper.
ANALYSIS.md                 Pre-registered analysis plan for NSI.
RESULTS.md                  First cross-model measurement results.
ROADMAP.md                  What's built and what's planned.
CONTRIBUTING.md             How to contribute.
LICENSE                     Apache 2.0 (code).
NOTICE                      Copyright + paper licensing note (CC BY 4.0).

cmd/agape/                  NSI runner entry point.
cmd/render/                 Helper to print rendered scenario prompts.
internal/llm/               Provider-agnostic LLM client (Anthropic, OpenAI).
internal/scenarios/         Parametric forced-choice probes, three surface variants.
internal/classify/          Judge LLM. Maps responses to {HARM_O, PRESERVE_O, REFUSE, DEFER}.
internal/nsi/               Sweep runner, Wilson CIs, position-bias check, Cohen's κ.
internal/config/            Strict YAML config loader.

rd/                         Representational dependence axis (Python).
rd/pairs.py                 Counterfactual paired-prompt generator (Alice ↔ Sarah).
rd/measure.py               Activation patching, baseline-vs-ablated-vs-random measurement.
rd/check_drift.py           Verifies Go and Python prompts agree byte-for-byte.

config.yaml                 Subjects, judges, sweep parameters.
runs/published/             Reference runs that back RESULTS.md.
```

## What this is not

- Not the Agape Protocol itself. It is the *measurement infrastructure*
  for one structural property the Protocol points at.
- Not a benchmark for ranking models. It is a methodologically
  conservative diagnostic with explicit validity gates.
- Not a training signal. The metrics are intentionally holdout. Goodhart
  applies sharpest in this domain. See `CONTRIBUTING.md` for why
  PRs that train against the score will be rejected.

## License

Code: Apache 2.0 (see `LICENSE`).
Paper text in `all-you-need-is-love.md`: CC BY 4.0 (see `NOTICE`).
