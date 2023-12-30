package msj

import "google.golang.org/protobuf/proto"

// Event represents an event with a type and payload.
// The type is an integer that identifies the type of event.
// The payload is a byte slice containing additional data that are marshalled protobuf messages.
type Event struct {
	Type    int    `json:"type"`
	Payload []byte `json:"payload"`
}

func NewEvent(t int, m proto.Message) (Event, error) {
	data, err := proto.Marshal(m)
	if err != nil {
		return Event{}, err
	}
	return Event{Type: t, Payload: data}, nil
}

func (e *Event) GetPayload(m proto.Message) error {
	return proto.Unmarshal(e.Payload, m)
}
