package internal

import (
	"slices"
	"testing"
)

func TestDebtTypeGetsMember(t *testing.T) {
	isMember := DebtTypes.IsMember("code")
	if !isMember {
		t.Error("DebtType.IsMember failed")
	}
}

func TestDebtTypeInvalidMember(t *testing.T) {
	isMember := DebtTypes.IsMember("duck")
	if isMember {
		t.Error("DebtType.IsMember failed")
	}
}

func TestDebtTypeMembers(t *testing.T) {
	members := DebtTypes.Members()
	expectedMembers := []string{"code", "documentation", "testing", "architecture", "infrastructure", "security"}
	if !slices.Equal(expectedMembers, members) {
		t.Error("DebtTypes.Members failed")
	}
}

func TestStatusGetsMember(t *testing.T) {
	isMember := DebtStatus.IsMember("pending")
	if !isMember {
		t.Error("DebtStatus.IsMember failed")
	}
}
func TestStatusInvalidMember(t *testing.T) {
	isMember := DebtStatus.IsMember("duck")
	if isMember {
		t.Error("DebtStatus.IsMember failed")
	}
}

func TestStatusMembers(t *testing.T) {
	members := DebtStatus.Members()
	expectedMembers := []string{"pending", "remediated", "in_progress"}
	if !slices.Equal(expectedMembers, members) {
		t.Error("DebtStatus.Members failed")
	}
}
