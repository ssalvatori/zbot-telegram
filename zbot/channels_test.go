package zbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelRegistryForget(t *testing.T) {
	r := newChannelRegistry([]Channel{{ID: 1, Title: "one"}, {ID: 2, Title: "two"}})

	r.Forget(1)
	assert.Equal(t, []Channel{{ID: 2, Title: "two"}}, r.Channels(), "first channel should be removed")

	r.Forget(2)
	assert.Empty(t, r.Channels(), "registry should be empty after forgetting all channels")
}

func TestChannelRegistryForgetUnknown(t *testing.T) {
	r := newChannelRegistry([]Channel{{ID: 1, Title: "one"}})

	r.Forget(999) // no-op, must not panic
	assert.Equal(t, []Channel{{ID: 1, Title: "one"}}, r.Channels(), "unknown ID should leave registry unchanged")
}

func TestChannelRegistryForgetFromEmpty(t *testing.T) {
	r := newChannelRegistry(nil)

	r.Forget(1) // no-op on empty registry
	assert.Empty(t, r.Channels())
}
