package httpapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCancelledCreateRequestDoesNotPersistVehicle(t *testing.T) {
	h := newAPIHarness(t)
	body := []byte(`{"plate_number":"沪A88888","vehicle_type":"sweeper","depot_code":"H-01","capacity_kg":6000,"odometer_km":10,"inspection_due_at":"2027-08-18T00:00:00Z"}`)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vehicles", bytes.NewReader(body)).WithContext(ctx)
	recorder := httptest.NewRecorder()
	h.handler.ServeHTTP(recorder, req)
	if recorder.Code == http.StatusCreated {
		t.Fatalf("cancelled request created vehicle: %s", recorder.Body.String())
	}
	page := h.request(t, http.MethodGet, "/api/v1/vehicles?limit=10&q=88888", nil, http.StatusOK)
	if page["total"].(float64) != 0 {
		t.Fatalf("cancelled request persisted data: %+v", page)
	}
}
