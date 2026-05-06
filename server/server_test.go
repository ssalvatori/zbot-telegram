package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChatID(t *testing.T) {

	channels := []Channel{
		{
			AuthToken: "authToken",
			ID:        1234,
		},
		{
			AuthToken: "authToken2",
			ID:        0,
		},
	}

	assert.Equal(t, int64(1234), getChatID("authToken", channels), "GET ID using token")
	assert.Equal(t, int64(0), getChatID("token2", channels), "GET ID using token")
}

func TestAPIMessages(t *testing.T) {
	channels := []Channel{
		{AuthToken: "valid-token", ID: 1234, Title: "test"},
	}

	handler := apiMessages(nil, channels)

	// No token → 403 Forbidden.
	req := httptest.NewRequest(http.MethodGet, "/messages", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Forbidden")

	// Invalid token (chatID == 0) → 403 Forbidden.
	req = httptest.NewRequest(http.MethodGet, "/messages?token=bad-token", nil)
	w = httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Valid token but no data → does not call bot.Send → 403 Forbidden.
	req = httptest.NewRequest(http.MethodGet, "/messages?token=valid-token", nil)
	w = httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
