package repositories

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		service     Service
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid service",
			service: Service{
				Name:        "TestService",
				ServiceType: "API",
				Description: "A test service",
				Url:         "https://test-service.com",
			},
			expectError: false,
		},
		{
			name: "Missing name",
			service: Service{
				ServiceType: "API",
				Description: "A test service",
				Url:         "https://test-service.com",
			},
			expectError: true,
			errorMsg:    "service name is required",
		},
		{
			name: "Missing URL",
			service: Service{
				Name:        "TestService",
				ServiceType: "API",
				Description: "A test service",
			},
			expectError: true,
			errorMsg:    "service url is required",
		},
		{
			name: "Missing service type",
			service: Service{
				Name:        "TestService",
				Description: "A test service",
				Url:         "https://test-service.com",
			},
			expectError: true,
			errorMsg:    "service type is required",
		},
		{
			name: "Invalid URL format",
			service: Service{
				Name:        "TestService",
				ServiceType: "API",
				Description: "A test service",
				Url:         "http://[invalid-url",
			},
			expectError: true,
			errorMsg:    "service url is not a valid URL format",
		},
		{
			name: "URL without scheme",
			service: Service{
				Name:        "TestService",
				ServiceType: "API",
				Description: "A test service",
				Url:         "invalid-url",
			},
			expectError: true,
			errorMsg:    "service url must use http or https protocol",
		},
		{
			name: "URL without http/https scheme",
			service: Service{
				Name:        "TestService",
				ServiceType: "API",
				Description: "A test service",
				Url:         "ftp://test-service.com",
			},
			expectError: true,
			errorMsg:    "service url must use http or https protocol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()

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
		})
	}
}
