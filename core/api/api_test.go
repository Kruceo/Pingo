package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"pingo/core/api"
	"pingo/core/config"
)

func TestItemsCRUDPersistsToFile(t *testing.T) {
	t.Parallel()

	cfgPath := filepath.Join(t.TempDir(), "config.json")
	if err := config.SaveConfig(cfgPath, &config.Config{
		PingInterval: 30,
		Items: []config.PingConfig{
			{
				Name:    "Cloudflare DNS IPv4",
				Tool:    "pingv4",
				Target:  "1.1.1.1",
				Timeout: 5000,
			},
		},
	}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	store := api.NewConfigStore(cfgPath)
	srv := httptest.NewServer(api.NewHandler(store))
	t.Cleanup(srv.Close)

	add := config.PingConfig{
		Name:    "Google DNS IPv4",
		Tool:    "pingv4",
		Target:  "8.8.8.8",
		Timeout: 3000,
	}
	resp := doJSON(t, srv.Client(), http.MethodPost, srv.URL+"/items", add)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, readAll(t, resp.Body))
	}
	_ = resp.Body.Close()

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(loaded.Items))
	}

	update := config.PingConfig{
		Name:    "Google DNS IPv4",
		Tool:    "pingv4",
		Target:  "8.8.8.8",
		Timeout: 1234,
	}
	itemURL := srv.URL + "/items/" + url.PathEscape(add.Name)
	resp = doJSON(t, srv.Client(), http.MethodPut, itemURL, update)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, readAll(t, resp.Body))
	}
	_ = resp.Body.Close()

	loaded, err = store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	found := false
	for _, it := range loaded.Items {
		if it.Name == add.Name {
			found = true
			if it.Timeout != 1234 {
				t.Fatalf("expected updated timeout 1234, got %d", it.Timeout)
			}
		}
	}
	if !found {
		t.Fatalf("updated item not found")
	}

	resp = doJSON(t, srv.Client(), http.MethodDelete, itemURL, nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", resp.StatusCode, readAll(t, resp.Body))
	}
	_ = resp.Body.Close()

	loaded, err = store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Items) != 1 {
		t.Fatalf("expected 1 item after delete, got %d", len(loaded.Items))
	}
}

func TestRejectsUnsupportedTool(t *testing.T) {
	t.Parallel()

	cfgPath := filepath.Join(t.TempDir(), "config.json")
	if err := config.SaveConfig(cfgPath, &config.Config{PingInterval: 30, Items: []config.PingConfig{}}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	store := api.NewConfigStore(cfgPath)
	srv := httptest.NewServer(api.NewHandler(store))
	t.Cleanup(srv.Close)

	resp := doJSON(t, srv.Client(), http.MethodPost, srv.URL+"/items", map[string]any{
		"name":    "Bad Tool",
		"tool":    "nope",
		"target":  "example.com",
		"timeout": 1000,
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", resp.StatusCode, readAll(t, resp.Body))
	}
	_ = resp.Body.Close()
}

func doJSON(t *testing.T, client *http.Client, method, url string, body any) *http.Response {
	t.Helper()

	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		r = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	return resp
}

func readAll(t *testing.T, r io.Reader) string {
	t.Helper()
	b, _ := io.ReadAll(r)
	return string(b)
}
