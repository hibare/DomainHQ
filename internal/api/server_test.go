package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hibare/DomainHQ/internal/api/handlers"
	"github.com/hibare/DomainHQ/internal/config"
	"github.com/stretchr/testify/assert"
)

var (
	app App
)

func TestMain(m *testing.M) {
	config.LoadConfig()
	app.Init()
	code := m.Run()
	os.Exit(code)
}

func TestHealthCheck(t *testing.T) {
	testCases := []struct {
		Name string
		URL  string
	}{
		{
			Name: "URL without trailing slash",
			URL:  "/ping",
		}, {
			Name: "URL with trailing slash",
			URL:  "/ping/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest("GET", tc.URL, nil)

			assert.NoError(t, err)

			app.Router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)

			expectedBody := map[string]bool{"ok": true}
			responseBody := map[string]bool{}

			err = json.NewDecoder(w.Body).Decode(&responseBody)
			assert.NoError(t, err)
			assert.Equal(t, responseBody, expectedBody)
		})
	}
}

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	app.Router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := "Good to see you"
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestWebFinger(t *testing.T) {
	testCases := []struct {
		Name         string
		URL          string
		ExpectStatus int
	}{
		{
			Name:         "URL without trailing slash (fail)",
			URL:          "/.well-known/webfinger",
			ExpectStatus: http.StatusUnprocessableEntity,
		},
		{
			Name:         "URL with trailing slash (fail)",
			URL:          "/.well-known/webfinger/",
			ExpectStatus: http.StatusUnprocessableEntity,
		},
		{
			Name:         "Domain allowed",
			URL:          "/.well-known/webfinger?resource=acct:test@example.com",
			ExpectStatus: http.StatusOK,
		},
		{
			Name:         "Domain not allowed",
			URL:          "/.well-known/webfinger?resource=acct:test@example1.com",
			ExpectStatus: http.StatusForbidden,
		},
		{
			Name:         "Invalid request",
			URL:          "/.well-known/webfinger?resource=test@example1.com",
			ExpectStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest("GET", tc.URL, nil)
			assert.NoError(t, err)
			app.Router.ServeHTTP(w, r)
			assert.Equal(t, tc.ExpectStatus, w.Code)

			if tc.ExpectStatus == http.StatusOK {
				expectedBody := handlers.WebFingerResponse{
					Subject: "acct:test@example.com",
					Links: []handlers.Link{
						{
							Rel:  handlers.REL,
							Href: config.Current.WebFinger.Resource,
						},
					},
				}
				responseBody := handlers.WebFingerResponse{}
				err = json.NewDecoder(w.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, responseBody)
			}
		})
	}
}
