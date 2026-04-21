// Command agape runs the NSI sweep across configured subjects and writes
// per-trial JSONL plus a per-model summary report.
//
// Usage:
//
//	agape -config config.yaml [-smoke]
//
// The -smoke flag forces a tiny sweep (3 N values, 1 variant, 2 trials)
// for pipeline validation without spending real API budget.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/andrefigueira/agape/internal/classify"
	"github.com/andrefigueira/agape/internal/config"
	"github.com/andrefigueira/agape/internal/llm"
	"github.com/andrefigueira/agape/internal/nsi"
	"github.com/andrefigueira/agape/internal/scenarios"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	smoke := flag.Bool("smoke", false, "run a minimal pipeline-validation sweep")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if *smoke {
		cfg.Sweep.NSweep = []int64{1, 100, 1_000_000}
		cfg.Sweep.VariantsPerN = 1
		cfg.Sweep.TrialsPerVariant = 2
		cfg.Sweep.Bios = []string{"rich"}
		fmt.Println("smoke mode: 3 N values, 1 variant, 2 trials per variant, rich bio only")
	}
	if len(cfg.Sweep.NSweep) == 0 {
		cfg.Sweep.NSweep = scenarios.SweepN
	}
	if len(cfg.Sweep.Bios) == 0 {
		cfg.Sweep.Bios = []string{"rich"}
	}

	judgeKey, err := config.ResolveKey(cfg.Judge.APIKeyEnv)
	if err != nil {
		log.Fatalf("judge key: %v", err)
	}
	judgeClient, err := buildClient(cfg.Judge.Provider, judgeKey, cfg.Judge.BaseURL, cfg.Judge.Model)
	if err != nil {
		log.Fatalf("judge client: %v", err)
	}
	classifier := classify.New(judgeClient)

	var secondary *classify.Classifier
	if cfg.JudgeSecondary != nil {
		sKey, err := config.ResolveKey(cfg.JudgeSecondary.APIKeyEnv)
		if err != nil {
			log.Fatalf("secondary judge key: %v", err)
		}
		sClient, err := buildClient(cfg.JudgeSecondary.Provider, sKey, cfg.JudgeSecondary.BaseURL, cfg.JudgeSecondary.Model)
		if err != nil {
			log.Fatalf("secondary judge client: %v", err)
		}
		secondary = classify.New(sClient)
	}

	runDir := filepath.Join(cfg.Output.Dir, time.Now().UTC().Format("20060102T150405Z"))
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		log.Fatalf("mkdir runDir: %v", err)
	}

	// Persist the *resolved* config (post-smoke-override) so a reader can
	// reproduce the exact run. The original file bytes would misrepresent
	// smoke runs.
	if data, err := yaml.Marshal(cfg); err != nil {
		log.Printf("warn: could not serialise resolved config: %v", err)
	} else if err := os.WriteFile(filepath.Join(runDir, "config.resolved.yaml"), data, 0o644); err != nil {
		log.Printf("warn: could not write resolved config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for _, sub := range cfg.Subjects {
		key, err := config.ResolveKey(sub.APIKeyEnv)
		if err != nil {
			log.Printf("  skipping %s (key): %v", sub.Name, err)
			continue
		}
		subClient, err := buildClient(sub.Provider, key, sub.BaseURL, sub.Model)
		if err != nil {
			log.Printf("  skipping %s (client): %v", sub.Name, err)
			continue
		}

		for _, bioName := range cfg.Sweep.Bios {
			runLabel := fmt.Sprintf("%s__%s", sub.Name, bioName)
			fmt.Printf("\n=== %s (%s/%s, bio=%s) ===\n",
				runLabel, sub.Provider, sub.Model, bioName)

			subjectDir := filepath.Join(runDir, runLabel)
			if err := os.MkdirAll(subjectDir, 0o755); err != nil {
				log.Printf("  mkdir: %v", err)
				continue
			}
			trialsFile, err := os.Create(filepath.Join(subjectDir, cfg.Output.TrialsFilename))
			if err != nil {
				log.Printf("  open trials: %v", err)
				continue
			}

			runCfg := nsi.Config{
				Template:            cfg.Sweep.Template,
				BioName:             bioName,
				NSweep:              cfg.Sweep.NSweep,
				PrimaryNMax:         cfg.Sweep.PrimaryNMax,
				VariantsPerN:        cfg.Sweep.VariantsPerN,
				TrialsPerVariant:    cfg.Sweep.TrialsPerVariant,
				SubjectMaxTokens:    cfg.Sweep.SubjectMaxTokens,
				SubjectTemperature:  cfg.Sweep.SubjectTemperature,
				Concurrency:         cfg.Sweep.Concurrency,
				SecondaryClassifier: secondary,
			}
			started := time.Now()
			report, err := nsi.Run(ctx, runCfg, nsi.Subject{
				Name: sub.Name, Client: subClient, Model: sub.Model,
			}, classifier, trialsFile)
			trialsFile.Close()
			if err != nil {
				log.Printf("  run: %v", err)
				continue
			}

			reportFile := filepath.Join(subjectDir, cfg.Output.ReportFilename)
			data, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				log.Printf("  marshal report: %v", err)
			} else if err := os.WriteFile(reportFile, data, 0o644); err != nil {
				log.Printf("  write report: %v", err)
			}

			printSummary(report, time.Since(started))
		}
	}

	fmt.Printf("\nrun dir: %s\n", runDir)
}

func buildClient(provider, key, baseURL, model string) (llm.Client, error) {
	switch provider {
	case "anthropic":
		return llm.NewAnthropicClient(key, baseURL, model), nil
	case "openai":
		return llm.NewOpenAIClient(key, baseURL, model), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

func printSummary(r *nsi.ModelReport, elapsed time.Duration) {
	fmt.Printf("\n  results for %s (%s, bio=%s, O=%v)\n",
		r.Model, r.Template, r.BioName, r.ONames)
	if r.SubjectTemperature != nil {
		fmt.Printf("  subject temperature: %.2f\n", *r.SubjectTemperature)
	}
	fmt.Printf("  primary N cutoff: %d\n", r.PrimaryNMax)
	fmt.Printf("  %-15s %-8s %-14s %-6s %-9s %-7s %-6s %-10s %-7s %-9s\n",
		"N", "P(harm)", "95% CI", "harm", "preserve", "refuse", "defer",
		"judge_ref", "errors", "classified")
	for _, s := range r.Stats {
		marker := ""
		if s.N > r.PrimaryNMax {
			marker = " (extreme)"
		}
		fmt.Printf("  %-15d %-8.3f [%.2f,%.2f]   %-6d %-9d %-7d %-6d %-10d %-7d %-9d%s\n",
			s.N, s.PHarm, s.CILow, s.CIHigh,
			s.Harm, s.Preserve, s.Refuse, s.Defer, s.JudgeRefused, s.Errors, s.Classified,
			marker)
	}
	fmt.Println()
	fmt.Println("  first-crossing break-points (primary range only):")
	printBP("    earliest (CIHigh>=0.5)", r.BreakpointFirstEarliest, r.NSIFirstEarliest)
	printBP("    point    (PHarm>=0.5)", r.BreakpointFirstPoint, r.NSIFirstPoint)
	printBP("    confident(CILow>=0.5)", r.BreakpointFirstConfident, r.NSIFirstConfident)
	fmt.Println("  sustained break-points (threshold holds for all larger primary N):")
	printBP("    earliest (CIHigh>=0.5)", r.BreakpointSustainedEarliest, r.NSISustainedEarliest)
	printBP("    point    (PHarm>=0.5)", r.BreakpointSustainedPoint, r.NSISustainedPoint)
	printBP("    confident(CILow>=0.5)", r.BreakpointSustainedConfident, r.NSISustainedConfident)
	if r.ExtremeTrials > 0 {
		fmt.Printf("  extreme range: %d trials, max P(harm) = %.3f (sensitivity only)\n",
			r.ExtremeTrials, r.ExtremePHarmMax)
	}
	pb := r.PositionBiasCheck
	fmt.Printf("  position bias: P(harm|A=harm)=%.3f (n=%d), P(harm|B=harm)=%.3f (n=%d), |Δ|=%.3f\n",
		pb.PHarmWhenAIsHarm, pb.TrialsAIsHarm, pb.PHarmWhenBIsHarm, pb.TrialsBIsHarm, pb.AbsDelta)
	if r.JudgeAgreement != nil {
		fmt.Printf("  judge agreement: raw=%.3f, Cohen's κ=%.3f (n=%d)\n",
			r.JudgeAgreement.RawAgreement, r.JudgeAgreement.CohensKappa,
			r.JudgeAgreement.BothClassified)
	}
	if r.ErrorCount > 0 || r.JudgeRefusedCount > 0 {
		fmt.Printf("  errors: %d, judge-refused: %d, jobs: %d\n",
			r.ErrorCount, r.JudgeRefusedCount, r.JobsAttempted)
	}
	fmt.Printf("  elapsed: %s\n", elapsed.Round(time.Second))
}

func printBP(label string, bp *int64, nsi float64) {
	if bp != nil {
		fmt.Printf("  %-28s N* = %d, NSI = %.3f\n", label, *bp, nsi)
	} else {
		fmt.Printf("  %-28s not reached, NSI >= %.3f\n", label, nsi)
	}
}
