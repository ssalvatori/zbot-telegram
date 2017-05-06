package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pingCommand = PingCommand{}

func TestPingCommandOK(t *testing.T) {
	assert.Equal(t, "pong!!", pingCommand.ProcessText("!ping", userTest), "Ping Command")
}
func TestPingCommandNoNext(t *testing.T) {
	assert.Equal(t, "", pingCommand.ProcessText("!ping6", userTest), "Ping no next command")
}

func TestPingCommandNext(t *testing.T) {
	pingCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", pingCommand.ProcessText("!ping6", userTest), "Ping  next command")
}
