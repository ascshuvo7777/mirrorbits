// Copyright (c) 2026 Amit Mishra
// Licensed under the MIT license

package rpc

import (
	"sort"
	"testing"
)

func names(mirrors []*MirrorID) []string {
	var out []string
	for _, m := range mirrors {
		out = append(out, m.Name)
	}
	sort.Strings(out)
	return out
}

func TestMatchMirrorsByPattern(t *testing.T) {
	// Regression test for https://github.com/videolabs/mirrorbits/issues/134
	tests := []struct {
		name    string
		mirrors map[int]string
		pattern string
		want    []string
	}{
		{
			name: "exact match takes priority over substring matches",
			mirrors: map[int]string{
				1: "fcix.net",
				2: "mirror.fcix.net",
				3: "paducahix.mm.fcix.net",
				4: "forksystems.mm.fcix.net",
			},
			pattern: "fcix.net",
			want:    []string{"fcix.net"},
		},
		{
			name: "exact match is case-insensitive",
			mirrors: map[int]string{
				1: "FCIX.net",
				2: "mirror.fcix.net",
			},
			pattern: "fcix.net",
			want:    []string{"FCIX.net"},
		},
		{
			name: "multiple substring matches returned when no exact match",
			mirrors: map[int]string{
				1: "mirror.fcix.net",
				2: "paducahix.mm.fcix.net",
			},
			pattern: "fcix.net",
			want:    []string{"mirror.fcix.net", "paducahix.mm.fcix.net"},
		},
		{
			name: "no match returns empty",
			mirrors: map[int]string{
				1: "alpha",
				2: "beta",
			},
			pattern: "gamma",
			want:    nil,
		},
		{
			name: "single substring match",
			mirrors: map[int]string{
				1: "mirror.example.com",
				2: "other.example.org",
			},
			pattern: "example.com",
			want:    []string{"mirror.example.com"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := names(matchMirrorsByPattern(tc.mirrors, tc.pattern))

			if len(got) != len(tc.want) {
				t.Fatalf("matchMirrorsByPattern(%q) = %v, want %v", tc.pattern, got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("matchMirrorsByPattern(%q) = %v, want %v", tc.pattern, got, tc.want)
				}
			}
		})
	}
}
