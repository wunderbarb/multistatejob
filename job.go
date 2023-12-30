// V0.0.1

package msj

import (
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/godruoyi/go-snowflake"
)

const (
	// JobPending represents the initial state of a job.
	JobPending JobState = iota + 1
	JobFailedAtStart
	JobFailed
	JobCompleted
	JobCustomSpace
)

// JobState represents the different states that a job can be in.
// The predefined states are:
//
// 1. JobPending: indicates that the job is pending and has not started yet.
// 2. JobFailedAtStart: indicates that the job could not start.
// 3. JobCompleted: indicates that the job has been completed successfully.
// 4. JobFailed: indicates that the job has failed to complete.
type JobState int32

// Job is the structure holding the information about a job.
type Job struct {
	Number    uint64    `json:"number"`
	Type      int64     `json:"type"`
	State     JobState  `json:"state"`
	Payload   []byte    `json:"payload"`
	Initiated time.Time `json:"initiated"`
	Ended     time.Time `json:"ended"`
}

// NewJob creates a new Job with the specified type and message payload.
// It marshals the message using protobuf and initializes the Job fields.
// The Number field is populated automatically.
// The State field is set to JobPending.
// The Initiated field is set to the current time using time.Now().
// The Payload field is set to the marshaled message.
// Returns the created Job or an error if marshaling fails.
func NewJob(typ int64, m proto.Message) (Job, error) {
	job := Job{
		Number:    snowflake.ID(),
		Type:      typ,
		State:     JobPending,
		Initiated: time.Now(),
	}
	err := job.setPayload(m)
	if err != nil {
		return Job{}, err
	}
	return job, nil
}

func (j *Job) GetPayload(m proto.Message) error {
	return proto.Unmarshal(j.Payload, m)
}

func (j *Job) Update(s JobState, m proto.Message) error {
	err := j.setPayload(m)
	if err != nil {
		return err
	}
	j.State = s
	return nil
}

func (j *Job) setPayload(m proto.Message) error {
	if m == nil {
		return nil
	}
	payload, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	j.Payload = payload
	return nil
}
