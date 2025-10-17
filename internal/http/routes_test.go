package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dilroop-us/ecommerce-go/internal/product"
)

func TestPing(t *testing.T) {
	r := Router(product.NewStore())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)
	if got := w.Body.String(); got != "pong" {
		t.Fatalf("want pong, got %q", got)
	}
}

func TestCreateAndListProducts(t *testing.T) {
	r := Router(product.NewStore())
	w := httptest.NewRecorder()

	body := bytes.NewBufferString(`{"name":"Coffee","price":9.99}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || len(w.Body.Bytes()) == 0 {
		t.Fatalf("list failed: %d %s", w.Code, w.Body.String())
	}
}
