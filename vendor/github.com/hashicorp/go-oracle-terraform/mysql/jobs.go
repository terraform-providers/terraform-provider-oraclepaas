package mysql

import (
	"fmt"
	"time"
)

// API URI Paths for the Root Job path
const (
	JobRootPath = "/paas/api/v1.1/activitylog/%s/job/%s"
)

// Default Poll Interval value
const waitForJobPollInterval = 1 * time.Second

// JobClient is a client for the Service functions of the Job API.
type JobClient struct {
	ResourceClient
	PollInterval time.Duration
	Timeout      time.Duration
}

// Jobs returns a JobClient for checking job status
func (c *MySQLClient) Jobs() *JobClient {
	return &JobClient{
		ResourceClient: ResourceClient{
			MySQLClient:      c,
			ResourceRootPath: JobRootPath,
		},
	}
}

// JobStatus defines the constants for the status of a job
type JobStatus string

const (
	// JobStatusNew - the job is new.
	JobStatusNew JobStatus = "NEW"
	// JobStatusRunning - the job is still running.
	JobStatusRunning JobStatus = "RUNNING"
	// JobStatusFailed - the job has failed.
	JobStatusFailed JobStatus = "FAILED"
	// JobStatusSucceed - the job has succeeded.
	JobStatusSucceed JobStatus = "SUCCEED"
)

// JobResponse details the job information received after submitting a request
type JobResponse struct {
	Details Details `json:"details"`
}

// Details details the attributes of the specific job that is running on the service instance
type Details struct {
	JobID   string `json:"jobId"`
	Message string `json:"message"`
}

// Job details the attributes related to a job
type Job struct {
	// Job ID
	ID int `json:"jobId"`
	// Status of the job
	Status JobStatus `json:"status"`
}

// GetJobInput specifies which job to retrieve
type GetJobInput struct {
	// ID of the job.
	// Required.
	ID string
}

// GetJob retrieves the job with the given id
func (c *JobClient) GetJob(getInput *GetJobInput) (*Job, error) {
	var job Job
	if err := c.getResource(getInput.ID, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// WaitForJobCompletion waits for a service instance to be in the desired state
func (c *JobClient) WaitForJobCompletion(input *GetJobInput, pollInterval, timeoutSeconds time.Duration) error {
	return c.client.WaitFor("job to complete", pollInterval, timeoutSeconds, func() (bool, error) {
		info, getErr := c.GetJob(input)
		if getErr != nil {
			return false, getErr
		}
		c.client.DebugLogString(fmt.Sprintf("job ID is %v", info.ID))
		switch s := info.Status; s {
		case JobStatusSucceed: // Target State
			c.client.DebugLogString("Job Succeeded")
			return true, nil
		case JobStatusFailed:
			c.client.DebugLogString("Job Failed")
			return false, fmt.Errorf("Job %q failed", input.ID)
		case JobStatusNew:
			c.client.DebugLogString("Job New")
			return false, nil
		case JobStatusRunning:
			c.client.DebugLogString("Job Running")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown job state: %s, waiting", s))
			return false, nil
		}
	})
}
