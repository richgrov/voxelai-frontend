package jobs

import (
	"errors"
	"time"
)

var ErrTimeout = errors.New("timed out")

type JobService interface {
	GetJob(id string) (*Job, error)
	StartJob(id string, prompt string) error
	UpdateJobStatus(id string, status Status, result string) error
	WaitForCompletion(id string, timeout time.Duration) (*Job, error)
}

type Job struct {
	Prompt string
	Result string
	Status Status
}

type Status int

const (
	StatusSucceess Status = iota
	StatusFailed
	StatusInProgress
)

func (status Status) String() string {
	switch status {
	case StatusSucceess:
		return "COMPLETED"
	case StatusFailed:
		return "FAILED"
	default:
		return ""
	}
}

func statusFromStr(status string) Status {
	switch status {
	case "COMPLETED":
		return StatusSucceess
	case "FAILED":
		return StatusFailed
	default:
		return StatusInProgress
	}
}
