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

func NewJob(typ int64, m proto.Message) (Job, error) {
	payload, err := proto.Marshal(m)
	if err != nil {
		return Job{}, err
	}
	job := Job{
		Number:    snowflake.ID(),
		Type:      typ,
		State:     JobPending,
		Payload:   payload,
		Initiated: time.Now(),
	}
	return job, nil
}

func (j *Job) GetPayload(m proto.Message) error {
	return proto.Unmarshal(j.Payload, m)
}
