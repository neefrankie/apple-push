package message

import (
	"github.com/sideshow/apns2/payload"
	"sync"
	"time"
)

// Message represents the data posted by web client APN CMS.
type Message struct {
	PageID string `json:"pageId"`
	Action string `json:"action"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Sound  bool   `json:"sound"`
	// Media is used to set `mutable-content`
	Media string `json:"media"`
	// ContentAvailable is used to set background update notification
	ContentAvailable bool `json:"contentAvailable"`
	// Determine the message is pushed to Apple development server or produciton server
	ProdServer   bool        `json:"prodServer"`
	CreatedBy    string      `json:"createdBy"`
	DeviceGroup  DeviceGroup `json:"deviceGroup"`
	ApnsID       string
	CreatedAt    time.Time
	DeviceCount  int
	InvalidCount int
	TimeElapsed  int
	mux          sync.Mutex
}

// Payload uses the Message to create a apns2 payload.Payload
func (m *Message) Payload() *payload.Payload {
	// Create a payload
	pl := payload.NewPayload()
	pl.Custom("id", m.PageID)
	pl.Custom("action", m.Action)

	// For background update notification, you should stop here. `apns` should be empty
	if m.ContentAvailable {
		pl.ContentAvailable()
		return pl
	}

	pl.AlertTitle(m.Title)
	pl.AlertBody(m.Body)
	pl.Badge(1)

	if m.Sound {
		pl.Sound("ping.aiff")
	}

	if m.Media != "" {
		pl.MutableContent()
		pl.Custom("media", m.Media)
	}

	return pl
}
