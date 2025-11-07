package internal

import (
	"github.com/google/uuid"
	"net/http"
	"time"
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

func GetDateFromRequestPath(varName string, req *http.Request) (time.Time, bool) {
	dateVal := req.PathValue(varName)
	date, err := time.Parse("2006-01-02", dateVal)
	return date, err == nil
}
