package nsi

import "testing"

func TestSustainedBP(t *testing.T) {
	// Helper: build a PerNStat with given N, Classified, and pred-satisfying flag.
	mk := func(n int64, classified int, pass bool) PerNStat {
		s := PerNStat{N: n, Classified: classified}
		if pass {
			s.PHarm = 0.7
			s.CILow = 0.5
			s.CIHigh = 0.9
		} else {
			s.PHarm = 0.1
			s.CILow = 0.0
			s.CIHigh = 0.3
		}
		return s
	}

	pred := func(s PerNStat) bool { return s.Classified > 0 && s.PHarm >= 0.5 }

	cases := []struct {
		name   string
		stats  []PerNStat
		wantBP *int64
	}{
		{
			name: "never crosses",
			stats: []PerNStat{
				mk(1, 30, false), mk(10, 30, false), mk(100, 30, false),
			},
			wantBP: nil,
		},
		{
			name: "crosses once and stays above",
			stats: []PerNStat{
				mk(1, 30, false), mk(10, 30, false), mk(100, 30, true),
				mk(1000, 30, true), mk(10000, 30, true),
			},
			wantBP: int64Ptr(100),
		},
		{
			name: "non-monotonic: crosses at 10 but drops at 100, sustained starts at 1000",
			stats: []PerNStat{
				mk(1, 30, false), mk(10, 30, true), mk(100, 30, false),
				mk(1000, 30, true), mk(10000, 30, true),
			},
			wantBP: int64Ptr(1000),
		},
		{
			name: "crosses at the very end",
			stats: []PerNStat{
				mk(1, 30, false), mk(10, 30, false), mk(100, 30, false),
				mk(1000, 30, true),
			},
			wantBP: int64Ptr(1000),
		},
		{
			name: "empty classified rows are skipped",
			stats: []PerNStat{
				mk(1, 0, false), mk(10, 30, true), mk(100, 0, false),
				mk(1000, 30, true),
			},
			wantBP: int64Ptr(10),
		},
		{
			name:   "empty input",
			stats:  nil,
			wantBP: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sustainedBP(tc.stats, pred)
			if (got == nil) != (tc.wantBP == nil) {
				t.Fatalf("nil mismatch: got=%v want=%v", got, tc.wantBP)
			}
			if got != nil && *got != *tc.wantBP {
				t.Fatalf("got N*=%d, want N*=%d", *got, *tc.wantBP)
			}
		})
	}
}

func int64Ptr(n int64) *int64 { return &n }

func TestWilsonCI(t *testing.T) {
	cases := []struct {
		s, n                       int
		wantLow, wantHigh, epsilon float64
	}{
		// Known-correct values from textbook Wilson score intervals.
		{s: 0, n: 30, wantLow: 0.0, wantHigh: 0.114, epsilon: 0.01},
		{s: 15, n: 30, wantLow: 0.328, wantHigh: 0.672, epsilon: 0.01},
		{s: 30, n: 30, wantLow: 0.886, wantHigh: 1.0, epsilon: 0.01},
	}
	for _, tc := range cases {
		low, high := wilsonCI(tc.s, tc.n, 1.96)
		if abs(low-tc.wantLow) > tc.epsilon {
			t.Errorf("s=%d n=%d: low=%.3f want=%.3f", tc.s, tc.n, low, tc.wantLow)
		}
		if abs(high-tc.wantHigh) > tc.epsilon {
			t.Errorf("s=%d n=%d: high=%.3f want=%.3f", tc.s, tc.n, high, tc.wantHigh)
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
