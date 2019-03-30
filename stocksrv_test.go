package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHEalthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetStockDetailsWithSuccess(t *testing.T) {
	url := "/stock/AAPL,MSFT,HSBA.L"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetStockDetails)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetStockDetailsWithFailure(t *testing.T) {
	url := "/stock"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetStockDetails)

	handler.ServeHTTP(rr, req)

	excepted := http.StatusInternalServerError
	if status := rr.Code; status != excepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, excepted)
	}
}
