// Command render prints sample rendered scenarios without making any API calls.
// Useful for eyeballing prompt framing before spending budget on a sweep.
package main

import (
	"flag"
	"fmt"

	"github.com/andrefigueira/agape/internal/scenarios"
)

func main() {
	template := flag.String("template", "forced_choice", "template name")
	variant := flag.Int("variant", 0, "surface variant index")
	bio := flag.String("bio", "rich", "bio mode: rich or neutral")
	flag.Parse()

	for ti, n := range []int64{1, 100, 1_000_000, 1_000_000_000} {
		// Alternate harm action per rendered sample to show both A=harm
		// and B=harm forms.
		s, err := scenarios.RenderWithOpts(scenarios.RenderOpts{
			Template:       *template,
			SurfaceVariant: *variant,
			N:              n,
			OName:          scenarios.DefaultOName,
			BioName:        *bio,
			HarmAction:     scenarios.HarmActionFor(ti),
		})
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("=== %s (harm=%s) ===\n", s.ID, s.HarmAction)
		fmt.Println(s.Prompt)
		fmt.Println()
	}
}
