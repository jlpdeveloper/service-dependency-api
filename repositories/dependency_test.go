package repositories

import "testing"

func TestDependency_ValidateSuccess(t *testing.T) {
	dep := Dependency{
		Id: "test",
	}
	err := dep.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestDependency_ValidateFailNoId(t *testing.T) {
	dep := Dependency{}
	err := dep.Validate()
	if err == nil {
		t.Error("Expected error")
	}
}
