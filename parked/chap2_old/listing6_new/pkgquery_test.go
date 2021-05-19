package pkgquery

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startTestPackageServer() *httptest.Server {
	pkgData := `[
{"name": "package1", "version": "1.1"},
{"name": "package2", "version": "1.0"}
]`
	mux := http.NewServeMux()
	mux.HandleFunc("/api/packages", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, pkgData)
	}))

	mux.HandleFunc("/packages", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "Packages available")
	}))

	ts := httptest.NewServer(mux)
	return ts
}

func TestFetchPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	client := createHTTPClientWithTimeout(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()
	req, err := createHTTPGetRequest(ctx, ts.URL+"/api/packages", nil)
	if err != nil {
		t.Fatal(err)
	}
	packages, err := fetchPackageData(client, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 2 {
		t.Logf("Expected 2 packages, Got back: %d", len(packages))
		t.Fail()
	}
	req, err = createHTTPGetRequest(ctx, ts.URL+"/packages", nil)
	if err != nil {
		t.Fatal(err)
	}
	packages, err = fetchPackageData(client, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 0 {
		t.Logf("Expected 0 packages, Got back: %d", len(packages))
		t.Fail()
	}
}