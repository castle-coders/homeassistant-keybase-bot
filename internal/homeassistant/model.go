package homeassistant

// HAEnvelope message envelope
type HAEnvelope struct {
	Type string `json:"type"`
}

// AuthMessage sent in response to auth_required
type AuthMessage struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

// SubscribeMessage to subscribe to specific event type
type SubscribeMessage struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

type errorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResultMessage expected as response to command in command phase
type ResultMessage struct {
	ID      int          `json:"id"`
	Type    string       `json:"type"`
	Success bool         `json:"success"`
	Error   errorMessage `json:"error"`
}

type eventData struct {
	Message string  `json:"message"`
	User    *string `json:"user,omitempty"`
	Team    *string `json:"team,omitempty"`
	Channel *string `json:"channel,omitempty"`
}

type eventInfo struct {
	EventType string    `json:"event_type"`
	Data      eventData `json:"data"`
}

// EventMessage for subscribed events
type EventMessage struct {
	ID    int       `json:"id"`
	Type  string    `json:"type"`
	Event eventInfo `json:"event"`
}
