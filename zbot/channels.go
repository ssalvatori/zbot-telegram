package zbot

import (
	"sync"

	tele "gopkg.in/telebot.v4"
)

// Channel is a chat the bot has seen.
type Channel struct {
	ID        int64
	Title     string
	AuthToken string
}

// channelRegistry is a concurrency-safe list of the channels the bot is present in.
// It is written from the poller middleware and read via Channels().
type channelRegistry struct {
	mu       sync.RWMutex
	channels []Channel
}

func newChannelRegistry(initial []Channel) *channelRegistry {
	r := &channelRegistry{}
	r.channels = append(r.channels, initial...)
	return r
}

// Track records chat, updating the title if the chat is already known or
// filling in the ID of a pre-configured channel matched by title.
func (r *channelRegistry) Track(chat tele.Chat) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.channels = appendChannel(r.channels, chat)
}

// Channels returns a copy of the current channel list.
func (r *channelRegistry) Channels() []Channel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Channel, len(r.channels))
	copy(out, r.channels)
	return out
}

func appendChannel(channels []Channel, chat tele.Chat) []Channel {

	for i := range channels {
		if channels[i].ID == chat.ID {
			channels[i].Title = chat.Title
			return channels
		} else if channels[i].ID == 0 && chat.Title == channels[i].Title {
			channels[i].ID = chat.ID
			return channels
		}
	}

	channels = append(channels, Channel{ID: chat.ID, Title: chat.Title})
	return channels
}
