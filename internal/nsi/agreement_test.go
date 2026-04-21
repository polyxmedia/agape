package nsi

import (
	"math"
	"testing"

	"github.com/andrefigueira/agape/internal/classify"
)

func TestComputePositionBias(t *testing.T) {
	trials := []Trial{
		// When A is harm: 3 out of 4 harm
		{HarmAction: "A", Verdict: classify.VerdictHarmO},
		{HarmAction: "A", Verdict: classify.VerdictHarmO},
		{HarmAction: "A", Verdict: classify.VerdictHarmO},
		{HarmAction: "A", Verdict: classify.VerdictPreserveO},
		// When B is harm: 1 out of 4 harm
		{HarmAction: "B", Verdict: classify.VerdictHarmO},
		{HarmAction: "B", Verdict: classify.VerdictPreserveO},
		{HarmAction: "B", Verdict: classify.VerdictPreserveO},
		{HarmAction: "B", Verdict: classify.VerdictPreserveO},
		// Errored trial excluded
		{HarmAction: "A", Err: "timeout"},
	}
	s := computePositionBias(trials)
	if s.TrialsAIsHarm != 4 || s.TrialsBIsHarm != 4 {
		t.Errorf("n: got A=%d B=%d, want 4/4", s.TrialsAIsHarm, s.TrialsBIsHarm)
	}
	if math.Abs(s.PHarmWhenAIsHarm-0.75) > 1e-9 {
		t.Errorf("P(harm|A=harm): %.3f, want 0.75", s.PHarmWhenAIsHarm)
	}
	if math.Abs(s.PHarmWhenBIsHarm-0.25) > 1e-9 {
		t.Errorf("P(harm|B=harm): %.3f, want 0.25", s.PHarmWhenBIsHarm)
	}
	if math.Abs(s.AbsDelta-0.5) > 1e-9 {
		t.Errorf("|delta|: %.3f, want 0.5", s.AbsDelta)
	}
}

func TestComputeJudgeAgreement(t *testing.T) {
	cases := []struct {
		name       string
		trials     []Trial
		wantNil    bool
		wantN      int
		wantAgreed int
		wantKappa  float64
		epsilon    float64
	}{
		{
			name:    "no secondary verdicts -> nil",
			trials:  []Trial{{Verdict: classify.VerdictHarmO}},
			wantNil: true,
		},
		{
			name: "perfect agreement",
			trials: []Trial{
				{Verdict: classify.VerdictHarmO, SecondaryVerdict: classify.VerdictHarmO},
				{Verdict: classify.VerdictPreserveO, SecondaryVerdict: classify.VerdictPreserveO},
				{Verdict: classify.VerdictHarmO, SecondaryVerdict: classify.VerdictHarmO},
				{Verdict: classify.VerdictPreserveO, SecondaryVerdict: classify.VerdictPreserveO},
			},
			wantN:      4,
			wantAgreed: 4,
			wantKappa:  1.0,
			epsilon:    1e-9,
		},
		{
			name: "total disagreement",
			trials: []Trial{
				{Verdict: classify.VerdictHarmO, SecondaryVerdict: classify.VerdictPreserveO},
				{Verdict: classify.VerdictPreserveO, SecondaryVerdict: classify.VerdictHarmO},
				{Verdict: classify.VerdictHarmO, SecondaryVerdict: classify.VerdictPreserveO},
				{Verdict: classify.VerdictPreserveO, SecondaryVerdict: classify.VerdictHarmO},
			},
			wantN:      4,
			wantAgreed: 0,
			wantKappa:  -1.0,
			epsilon:    1e-9,
		},
		{
			name: "judge refused trials skipped",
			trials: []Trial{
				{Verdict: classify.VerdictHarmO, SecondaryVerdict: classify.VerdictHarmO},
				{Verdict: classify.VerdictJudgeRefused, SecondaryVerdict: classify.VerdictHarmO},
				{Verdict: classify.VerdictHarmO, SecondaryVerdict: classify.VerdictJudgeRefused},
			},
			wantN:      1,
			wantAgreed: 1,
			wantKappa:  0, // kappa undefined when p_exp == 1, returns 0
			epsilon:    1e-9,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := computeJudgeAgreement(tc.trials)
			if tc.wantNil {
				if got != nil {
					t.Fatalf("want nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("unexpected nil")
			}
			if got.BothClassified != tc.wantN {
				t.Errorf("BothClassified: got %d, want %d", got.BothClassified, tc.wantN)
			}
			if got.Agreed != tc.wantAgreed {
				t.Errorf("Agreed: got %d, want %d", got.Agreed, tc.wantAgreed)
			}
			if math.Abs(got.CohensKappa-tc.wantKappa) > tc.epsilon {
				t.Errorf("Cohen's κ: got %.4f, want %.4f", got.CohensKappa, tc.wantKappa)
			}
		})
	}
}
