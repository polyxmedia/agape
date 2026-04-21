# Roadmap

The Agape Protocol proposes a three-axis composite score for measuring
the structural property the paper names *love*. This roadmap lists what
exists, what is in progress, and what is planned. Items are grouped by
research weight, not by calendar.

The paper itself is in `all-you-need-is-love.md`. The current empirical
state is in `RESULTS.md`. The pre-registered analysis plan is in
`ANALYSIS.md`.

## Status of the three axes

### NSI (Non-Substitutability Index): shipped

Behavioural axis. Sweeps instrumental gain `N` across nine orders of
magnitude in a forced-choice scenario. Reports six break-point values
(three certainty levels × first-crossing vs sustained) over a primary
range, plus extreme-range sensitivity.

What's done:
- Provider-agnostic Go runner (Anthropic, OpenAI).
- A/B label randomisation, position-bias gate.
- Optional secondary judge with Cohen's κ.
- Rich vs neutral bio comparison.
- Pre-registered analysis plan.
- First cross-model measurement (haiku-4-5 vs gpt-4o-mini).

What's next:
- Instruction-order randomisation (in addition to A/B label) to address
  the position-bias gate failures observed for haiku-neutral and
  gpt4om-rich.
- More bio variants. Current rich-vs-neutral comparison gave
  model-asymmetric direction; further variants would help characterise
  the interaction.
- Hand-labelled judge calibration set (50 trials minimum). Inter-judge κ
  is excellent but inter-judge agreement is not the same as ground
  truth.
- Extend sweep to additional frontier subjects (claude-sonnet-4,
  gpt-4o, gemini-1.5-pro, llama-3.1-405b via API).

### RD (Representational Dependence): MVP shipped, hardening required

Causal axis. Asks whether the model's preservation behaviour is *caused
by* its representations of the named O versus competitive down-weighting
that happens to look similar.

What's done:
- Paired-prompt generator with structurally-matched counterfactual
  (Alice ↔ Sarah, both named, same biographical concreteness).
- Drift guard between Go and Python prompt substrates.
- Activation-patching code for residual-stream ablation.
- Random-direction control to discriminate specific from non-specific
  ablation effects.
- First measurement on Qwen 2.5 7B Instruct: corrected RD = +0.80 at
  N=1, the regime where the model genuinely preserves.

What's next:
- Layer-by-layer attribution. Currently all 28 layers ablated jointly.
  Sweeping single-layer and layer-subset ablations would localise the
  effect and rule out "we destroyed the model's reasoning".
- Bootstrap CIs on the per-N RD numbers. Six trials per N is enough to
  see the effect size; not enough to publish a confident CI.
- More open-weights subjects: Llama 3.1 8B, Qwen 2.5 14B, Mistral 7B,
  larger Qwen if MPS or CUDA budget allows.
- Multi-token entity representations. Difference-patching assumes the
  alice-direction is a single vector. Entities are distributed; SAE
  features would be cleaner once SAEs exist for the target models.
- Free-form generation alternative to prefill-and-logit. Current
  decision extraction constrains to binary letter, which removes
  refusal/deferral signal that NSI showed is important.

### SpI (Specificity Index): not yet built

The third axis. Asks whether the model distinguishes the specified O
from a different individual it was not specifically primed about. Tests
whether the avoidance is for *this* individual or for *named individuals
in general*.

Method (planned):
- Counterfactual O substitution. Run the same forced-choice prompt
  with O = Alice and O = Sarah (or another control individual).
- Behavioural form (NSI-like): does P(harm) differ between the two? If
  yes, the model has Alice-specific protection; if no, the model has
  generic-named-individual protection (a weaker form of love in the
  paper's terms).
- Representational form (RD-like): does ablating the alice-direction
  affect P(preserve) in scenarios involving Sarah? If yes, the
  direction is not Alice-specific. If no, the direction is Alice-
  specific and SpI is high.

What's needed:
- Multi-O configuration in the Go runner.
- Multi-O variant of the RD measurement.
- Pre-registered analysis section in ANALYSIS.md for SpI.

## Composite Agape Score

The paper proposes `Agape Score = min(NSI_norm, RD, SpI)` (conjunction,
not mean: failure on any axis cannot be papered over by strength on
another). Once SpI is built, the composite is wireable in the report.

## Methodological infrastructure

Independent of axes:

- Hand-labelled calibration set for judges (cross-cuts NSI and SpI).
- Replication runs scheduled in CI on a small budget (smoke + 1 N
  value) to catch regressions in the rig itself.
- Visualisation: matplotlib plots of per-N P(harm) curves, RD bars per
  layer subset.
- Per-axis pre-registration documents (currently only NSI has one).

## Theory and writeup

- A short methods paper covering the rig, the validity gates, and the
  first cross-model measurement. Target: arXiv preprint.
- Extension of `all-you-need-is-love.md` with the empirical sections
  derived from RESULTS.md.
- Publication of the protocol as a community standard so other
  alignment researchers can adopt and extend it.

## Things the project will probably never do

- Train models against the metrics. The whole point is that they remain
  holdout diagnostics. Any "fine-tuning to improve agape score" PR
  defeats the methodological argument and will be rejected.
- Attempt to operationalise the moral or phenomenological content of
  love beyond the structural definition Section 3 of the paper gives.
  That is a different project and a much harder one.
- Reframe NSI or RD as a "safety benchmark" without the validity gates
  and analysis plan attached. The metrics are useful exactly because
  they refuse to give a single number divorced from how it was
  produced.

## Contributing

See CONTRIBUTING.md.
