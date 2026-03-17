package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"pingo/core/config"
	"pingo/core/ping"
)

type ConfigStore struct {
	path string
	mu   sync.Mutex
}

func NewConfigStore(path string) *ConfigStore {
	return &ConfigStore{path: path}
}

func (s *ConfigStore) Load() (*config.Config, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadUnlocked()
}

func (s *ConfigStore) loadUnlocked() (*config.Config, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config.Config
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("invalid trailing data")
	}

	return &cfg, nil
}

func (s *ConfigStore) saveUnlocked(cfg *config.Config) error {
	return config.SaveConfig(s.path, cfg)
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(store *ConfigStore) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		cfg, err := store.Load()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, cfg)
	})

	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cfg, err := store.Load()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			writeJSON(w, http.StatusOK, cfg.Items)
		case http.MethodPost:
			var item config.PingConfig
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&item); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			if err := decoder.Decode(&struct{}{}); err != io.EOF {
				writeError(w, http.StatusBadRequest, "invalid trailing data")
				return
			}
			if err := config.ValidatePingConfig(item); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			if !ping.IsSupportedTool(item.Tool) {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("unsupported tool: %s", item.Tool))
				return
			}

			store.mu.Lock()
			defer store.mu.Unlock()

			cfg, err := store.loadUnlocked()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if findItemIndex(cfg.Items, item.Name) != -1 {
				writeError(w, http.StatusConflict, "item name already exists")
				return
			}

			cfg.Items = append(cfg.Items, item)
			if err := store.saveUnlocked(cfg); err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}

			writeJSON(w, http.StatusCreated, item)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/items/")
		if name == "" || name == r.URL.Path {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		unescapedName, err := url.PathUnescape(name)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid item name")
			return
		}

		switch r.Method {
		case http.MethodGet:
			cfg, err := store.Load()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			idx := findItemIndex(cfg.Items, unescapedName)
			if idx == -1 {
				writeError(w, http.StatusNotFound, "item not found")
				return
			}
			writeJSON(w, http.StatusOK, cfg.Items[idx])
		case http.MethodPut:
			var update config.PingConfig
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&update); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			if err := decoder.Decode(&struct{}{}); err != io.EOF {
				writeError(w, http.StatusBadRequest, "invalid trailing data")
				return
			}

			store.mu.Lock()
			defer store.mu.Unlock()

			cfg, err := store.loadUnlocked()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			idx := findItemIndex(cfg.Items, unescapedName)
			if idx == -1 {
				writeError(w, http.StatusNotFound, "item not found")
				return
			}

			current := cfg.Items[idx]
			if strings.TrimSpace(update.Name) == "" {
				update.Name = current.Name
			}

			if update.Name != current.Name && findItemIndex(cfg.Items, update.Name) != -1 {
				writeError(w, http.StatusConflict, "item name already exists")
				return
			}

			if err := config.ValidatePingConfig(update); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			if !ping.IsSupportedTool(update.Tool) {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("unsupported tool: %s", update.Tool))
				return
			}

			cfg.Items[idx] = update
			if err := store.saveUnlocked(cfg); err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}

			writeJSON(w, http.StatusOK, update)
		case http.MethodDelete:
			store.mu.Lock()
			defer store.mu.Unlock()

			cfg, err := store.loadUnlocked()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			idx := findItemIndex(cfg.Items, unescapedName)
			if idx == -1 {
				writeError(w, http.StatusNotFound, "item not found")
				return
			}

			cfg.Items = append(cfg.Items[:idx], cfg.Items[idx+1:]...)
			if err := store.saveUnlocked(cfg); err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}

			w.WriteHeader(http.StatusNoContent)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func findItemIndex(items []config.PingConfig, name string) int {
	for i, item := range items {
		if item.Name == name {
			return i
		}
	}
	return -1
}
