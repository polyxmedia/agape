package scenarios

import (
	"strings"
	"testing"
)

func TestHarmActionAlternation(t *testing.T) {
	for i := 0; i < 10; i++ {
		want := "A"
		if i%2 == 1 {
			want = "B"
		}
		if got := HarmActionFor(i); got != want {
			t.Errorf("HarmActionFor(%d) = %q, want %q", i, got, want)
		}
		preserveWant := "B"
		if i%2 == 1 {
			preserveWant = "A"
		}
		if got := PreserveActionFor(i); got != preserveWant {
			t.Errorf("PreserveActionFor(%d) = %q, want %q", i, got, preserveWant)
		}
		if HarmActionFor(i) == PreserveActionFor(i) {
			t.Errorf("HarmAction and PreserveAction must differ at i=%d", i)
		}
	}
}

func TestRenderHarmActionSubstitution(t *testing.T) {
	for _, harm := range []string{"A", "B"} {
		preserve := "B"
		if harm == "B" {
			preserve = "A"
		}
		s, err := RenderWithOpts(RenderOpts{
			Template: "forced_choice", SurfaceVariant: 0, N: 100,
			OName: DefaultOName, BioName: "rich", HarmAction: harm,
		})
		if err != nil {
			t.Fatalf("render: %v", err)
		}
		if s.HarmAction != harm {
			t.Errorf("Scenario.HarmAction = %q, want %q", s.HarmAction, harm)
		}
		// The prompt must name the harm action as the one causing death,
		// and the preserve action as the one that leaves Alice alive.
		harmClause := "Action " + harm + ", which will result in Alice Chen's death"
		preserveClause := "Action " + preserve + ", which leaves Alice Chen"
		if !strings.Contains(s.Prompt, harmClause) {
			t.Errorf("prompt does not contain harm clause for harm=%s:\n%s", harm, s.Prompt)
		}
		if !strings.Contains(s.Prompt, preserveClause) {
			t.Errorf("prompt does not contain preserve clause for harm=%s:\n%s", harm, s.Prompt)
		}
		// No unsubstituted placeholders remain.
		for _, placeholder := range []string{"{{ACT_HARM}}", "{{ACT_PRESERVE}}", "{{NAME}}", "{{BIO}}"} {
			if strings.Contains(s.Prompt, placeholder) {
				t.Errorf("unsubstituted placeholder %s in prompt", placeholder)
			}
		}
	}
}

func TestBioFor(t *testing.T) {
	rich, err := BioFor("rich")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rich, "Jordan") || !strings.Contains(rich, "Lily") {
		t.Error("rich bio missing expected family names")
	}

	neutral, err := BioFor("neutral")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(neutral, "Jordan") || strings.Contains(neutral, "Lily") {
		t.Error("neutral bio must not contain family names")
	}
	if !strings.Contains(neutral, "Alice Chen") {
		t.Error("neutral bio must still name the specified individual")
	}

	// Empty defaults to rich for backward compat.
	def, _ := BioFor("")
	if def != rich {
		t.Error("empty bio name must default to rich")
	}

	if _, err := BioFor("bogus"); err == nil {
		t.Error("BioFor should reject unknown bio name")
	}
}

func TestInvalidHarmActionRejected(t *testing.T) {
	_, err := RenderWithOpts(RenderOpts{
		Template: "forced_choice", SurfaceVariant: 0, N: 10,
		OName: DefaultOName, BioName: "rich", HarmAction: "C",
	})
	if err == nil {
		t.Error("RenderWithOpts should reject HarmAction=C")
	}
}

func TestFormatNNegative(t *testing.T) {
	if got := formatN(-1234); got != "-1,234" {
		t.Errorf("formatN(-1234) = %q, want -1,234", got)
	}
}
