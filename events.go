package msj

import "google.golang.org/protobuf/proto"

const (

	// EventCompleted is a constant representing the type of the event signaling a job successful completion.
	EventCompleted = 1 + iota
	EventFailed
	EventCustomSpace
)

// Event represents an event with a type and payload.
// The type is an integer that identifies the type of event.
// The payload is a byte slice containing additional data that are marshalled protobuf messages.
type Event struct {
	Job     uint64 `json:"job"`
	Type    int    `json:"type"`
	Payload []byte `json:"payload"`
}

func NewEvent(jn uint64, t int, m proto.Message) (Event, error) {
	data, err := proto.Marshal(m)
	if err != nil {
		return Event{}, err
	}
	return Event{Type: t, Payload: data, Job: jn}, nil
}

func (e *Event) GetPayload(m proto.Message) error {
	return proto.Unmarshal(e.Payload, m)
}
