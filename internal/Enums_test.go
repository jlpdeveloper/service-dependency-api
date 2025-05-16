package internal

import (
	"slices"
	"testing"
)

func TestEnumMembership(t *testing.T) {
	tests := []struct {
		name     string
		enum     StringEnum
		value    string
		expected bool
	}{
		{"ValidDebtType", DebtTypes, "code", true},
		{"InvalidDebtType", DebtTypes, "duck", false},
		{"ValidDebtStatus", DebtStatus, "pending", true},
		{"InvalidDebtStatus", DebtStatus, "duck", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.enum.IsMember(tc.value)
			if result != tc.expected {
				t.Errorf("%s: IsMember(%q) = %v, want %v", tc.name, tc.value, result, tc.expected)
			}
		})
	}
}

func TestEnumMembers(t *testing.T) {
	tests := []struct {
		name     string
		enum     StringEnum
		expected []string
	}{
		{"DebtTypes", DebtTypes, []string{"code", "documentation", "testing", "architecture", "infrastructure", "security"}},
		{"DebtStatus", DebtStatus, []string{"pending", "remediated", "in_progress"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			members := tc.enum.Members()
			if !slices.Equal(members, tc.expected) {
				t.Errorf("%s: Members() = %v, want %v", tc.name, members, tc.expected)
			}
		})
	}
}
