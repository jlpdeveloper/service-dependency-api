package internal

import (
	"github.com/google/uuid"
	"net/http"
)

type PathValidator func(string, *http.Request) (string, bool)

func GetGuidFromRequestPath(varName string, req *http.Request) (string, bool) {
	guidVal := req.PathValue(varName)
	return IsValidGuid(guidVal)
}

func IsValidGuid(guidVal string) (string, bool) {
	err := uuid.Validate(guidVal)
	return guidVal, err == nil
}
