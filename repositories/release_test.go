package repositories

import (
	"encoding/json"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestRelease_Validate(t *testing.T) {
	validUUID := uuid.New().String()
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		release     Release
		expectError bool
		errorMsg    string
		checkTime   bool
	}{
		{
			name: "Valid release with valid UUID",
			release: Release{
				ServiceId:   validUUID,
				ReleaseDate: fixedTime,
				Url:         "https://example.com",
			},
			expectError: false,
			checkTime:   false,
		},
		{
			name: "Invalid UUID",
			release: Release{
				ServiceId:   "not-a-uuid",
				ReleaseDate: fixedTime,
				Url:         "https://example.com",
			},
			expectError: true,
			errorMsg:    "invalid Service Id",
			checkTime:   false,
		},
		{
			name: "Zero release date should be set to current time",
			release: Release{
				ServiceId:   validUUID,
				ReleaseDate: time.Time{}, // Zero time
				Url:         "https://example.com",
			},
			expectError: false,
			checkTime:   true,
		},
		{
			name: "Non-zero release date should remain unchanged",
			release: Release{
				ServiceId:   validUUID,
				ReleaseDate: fixedTime,
				Url:         "https://example.com",
			},
			expectError: false,
			checkTime:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the release to avoid modifying the original in the test table
			release := tt.release

			// Store the original release date if we need to check it later
			originalDate := release.ReleaseDate

			// Call the Validate method
			err := release.Validate()

			// Check if we got the expected error
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}

			// Check if the release date was modified as expected
			if tt.checkTime {
				if release.ReleaseDate.IsZero() {
					t.Errorf("Expected release date to be set to current time, but it's still zero")
				}

				// Check that the time is recent (within the last minute)
				now := time.Now().UTC()
				diff := now.Sub(release.ReleaseDate)
				if diff < 0 || diff > time.Minute {
					t.Errorf("Expected release date to be close to current time, but got: %v (diff: %v)", release.ReleaseDate, diff)
				}
			} else if !tt.expectError && !originalDate.IsZero() {
				// For non-zero dates that should remain unchanged
				if !release.ReleaseDate.Equal(originalDate) {
					t.Errorf("Expected release date to remain %v, but got: %v", originalDate, release.ReleaseDate)
				}
			}
		})
	}
}

func TestRelease_FromJson(t *testing.T) {
	j := `{
  	"service_id": "12973952-0165-400b-9178-a0fdbd90f967",
	"url": "http://example.com"
	}`
	r := &Release{}
	err := json.Unmarshal([]byte(j), r)
	if err != nil {
		t.Errorf("Error unmarshalling json: %v", err)
	}
	err = r.Validate()
	if err != nil {
		t.Errorf("Error validating json: %v", err)
	}
	if r.ReleaseDate.IsZero() {
		t.Errorf("Expected release date to be set to current time, but it's still zero")
	}

	// Check that the time is recent (within the last minute)
	now := time.Now().UTC()
	diff := now.Sub(r.ReleaseDate)
	if diff < 0 || diff > time.Minute {
		t.Errorf("Expected release date to be close to current time, but got: %v (diff: %v)", r.ReleaseDate, diff)
	}
}
