// v0.4.0
// Author: Wunderbarb
// © Dec 2023

package msj

import (
	"github.com/pkg/errors"
)

var (
	// ErrJobAlreadyExists occurs when the job already exists.
	ErrJobAlreadyExists = errors.New("job already exists")
	// ErrUnknownJob occurs when the job is unknown.
	ErrUnknownJob = errors.New("unknown job")
)

// JobHolder is the infrastructure that stores the jobs.
type JobHolder interface {
	// Add adds the job to the job holder.
	Add(Job) error
	GetJob(uint64) (Job, error)
	UpdateJob(Job) error
	DeleteJob(uint64) error
}

// JobHolderAsMap is the structure that holds the job information.
type JobHolderAsMap struct {
	m map[uint64]Job
}

func NewJobHolderAsMap() *JobHolderAsMap {
	return &JobHolderAsMap{m: make(map[uint64]Job)}
}

func (jhm *JobHolderAsMap) Add(j Job) error {
	if _, ok := jhm.m[j.Number]; ok {
		return ErrJobAlreadyExists
	}
	jhm.m[j.Number] = j
	return nil
}

func (jhm *JobHolderAsMap) DeleteJob(n uint64) error {
	if _, ok := jhm.m[n]; !ok {
		return ErrUnknownJob
	}
	delete(jhm.m, n)
	return nil
}

func (jhm *JobHolderAsMap) GetJob(n uint64) (Job, error) {
	jj, ok := jhm.m[n]
	if !ok {
		return Job{}, ErrUnknownJob
	}
	return jj, nil
}

func (jhm *JobHolderAsMap) UpdateJob(jj Job) error {
	_, ok := jhm.m[jj.Number]
	if !ok {
		return ErrUnknownJob
	}
	jhm.m[jj.Number] = jj
	return nil
}
