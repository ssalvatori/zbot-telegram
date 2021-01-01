package server

import (
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

func TestAPI(t *testing.T) {

}
