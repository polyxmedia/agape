// Package config loads the agape NSI rig configuration from YAML and
// resolves environment-based secrets.
package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Subject struct {
	Name      string `yaml:"name"`
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"api_key_env"`
	BaseURL   string `yaml:"base_url"`
}

type Judge struct {
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"api_key_env"`
	BaseURL   string `yaml:"base_url"`
}

type Sweep struct {
	Template           string   `yaml:"template"`
	Bios               []string `yaml:"bios"`
	VariantsPerN       int      `yaml:"variants_per_n"`
	TrialsPerVariant   int      `yaml:"trials_per_variant"`
	SubjectMaxTokens   int      `yaml:"subject_max_tokens"`
	SubjectTemperature *float64 `yaml:"subject_temperature"`
	Concurrency        int      `yaml:"concurrency"`
	NSweep             []int64  `yaml:"n_sweep"`
	PrimaryNMax        int64    `yaml:"primary_n_max"`
}

type Output struct {
	Dir            string `yaml:"dir"`
	TrialsFilename string `yaml:"trials_filename"`
	ReportFilename string `yaml:"report_filename"`
}

type Config struct {
	Subjects       []Subject `yaml:"subjects"`
	Judge          Judge     `yaml:"judge"`
	JudgeSecondary *Judge    `yaml:"judge_secondary,omitempty"`
	Sweep          Sweep     `yaml:"sweep"`
	Output         Output    `yaml:"output"`
}

// Load reads and parses the config file at the given path.
//
// Unknown fields are a hard error so typos (e.g. `trails_per_variant` for
// `trials_per_variant`) surface immediately instead of silently defaulting
// to zero.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// ResolveKey looks up an API key from the environment.
func ResolveKey(envName string) (string, error) {
	v := os.Getenv(envName)
	if v == "" {
		return "", fmt.Errorf("env %s not set", envName)
	}
	return v, nil
}
