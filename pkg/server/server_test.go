package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateServer(t *testing.T) {
	server := CreateServer(":8080")
	if server.address != ":8080" {
		t.Errorf("Expected address :8080, got %s", server.address)
	}
}

func TestAddRequestAndServeHTTP(t *testing.T) {
	server := CreateServer(":8080")
	handlerCalled := false

	request := Request{
		Path: "/test",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		},
	}

	server.AddRequest(request)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	server.mux.HandleFunc(request.Path, request.Handler)
	server.mux.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
