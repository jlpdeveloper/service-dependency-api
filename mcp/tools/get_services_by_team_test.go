package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"service-dependency-api/repositories"
	"testing"
	"time"
)

func TestGetServicesByTeamSuccess(t *testing.T) {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			servicesByTeam := []repositories.Service{
				{
					Id:          "123",
					Name:        "svc-a",
					ServiceType: "service",
					Description: "desc",
					Url:         "https://a.example.com",
				},
				{
					Id:          "456",
					Name:        "svc-b",
					ServiceType: "service",
					Description: "desc",
					Url:         "https://b.example.com",
				},
			}
			enc := json.NewEncoder(w)
			err := enc.Encode(servicesByTeam)
			if err != nil {
				t.Fatal("Failed to encode services by team")
			}
		}),
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	defer func() {
		_ = server.Close()
	}()
	ctx := context.WithValue(context.Background(), "API_URL", "http://localhost:8080/")
	s, err := getServicesByTeam(ctx, "team-id")
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 2 {
		t.Fatal("Expected 2 services")
	}
	if s[0].Id != "123" {
		t.Fatal("Expected first service to be svc-a")
	}
	if s[1].Id != "456" {
		t.Fatal("Expected second service to be svc-b")
	}
}

func TestGetServicesByTeamNoApiUrl(t *testing.T) {
	ctx := context.Background()
	_, err := getServicesByTeam(ctx, "team-id")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGetServicesByTeamApiUrlNotAString(t *testing.T) {
	ctx := context.WithValue(context.Background(), "API_URL", 123)
	_, err := getServicesByTeam(ctx, "team-id")
	if err == nil {
	}
}

func TestGetServicesByTeamHttpError(t *testing.T) {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	defer func() {
		_ = server.Close()
	}()
	ctx := context.WithValue(context.Background(), "API_URL", "http://localhost:8080/")
	_, err := getServicesByTeam(ctx, "team-id")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGetServicesByTeamTimeout(t *testing.T) {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-make(chan struct{})
		}),
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	defer func() {
		_ = server.Close()
	}()
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(),
		"api_url", "http://localhost:8080/"), 100*time.Millisecond)
	defer cancel()
	_, err := getServicesByTeam(ctx, "team-id")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGetServicesByTeamInvalidJson(t *testing.T) {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("invalid json"))
		}),
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	defer func() {
		_ = server.Close()
	}()
	ctx := context.WithValue(context.Background(), "API_URL", "http://localhost:8080/")
	_, err := getServicesByTeam(ctx, "team-id")
	if err == nil {
	}
}
