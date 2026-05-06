package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pingCommand = PingCommand{}

func TestPingCommandOK(t *testing.T) {
	result, _ := pingCommand.ProcessText("!ping", userTest, "testchat", false)
	assert.Equal(t, "pong!!", result, "Ping Command")
}

func TestPingCommandNotMatch(t *testing.T) {

	result, err := pingCommand.ProcessText("!ping6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")
	assert.Equal(t, ErrNextCommand, err, "Error output doesn't match")
}
