package repositories

import "testing"

func TestDebt_ValidateSuccess(t *testing.T) {
	debt := Debt{
		Title:     "test",
		Type:      "code",
		ServiceId: "test",
	}
	err := debt.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestDebt_ValidateFailNoTitle(t *testing.T) {
	debt := Debt{
		Type:      "code",
		ServiceId: "test",
	}
	err := debt.Validate()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestDebt_ValidateFailNoType(t *testing.T) {
	debt := Debt{
		Title:     "test",
		ServiceId: "test",
	}
	err := debt.Validate()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestDebt_ValidateFailIncorrectType(t *testing.T) {
	debt := Debt{
		Title:     "test",
		Type:      "duck",
		ServiceId: "test",
	}
	err := debt.Validate()
	if err == nil {
		t.Error("Expected error")
	}
}
