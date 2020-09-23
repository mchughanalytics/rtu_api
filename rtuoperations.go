package rtu_api

import (
	"time"
)

/*
	RTU OPERATIONS

	Valid states: QUEUED, RUNNING, DONE, CANCELED
*/

// RtuOperation
type RtuOperation struct {
	Name          string
	ID            string
	created       time.Time
	started       time.Time
	completed     time.Time
	operation     func(client *RmClient, data DataSet) (i interface{}, e error)
	client        *RmClient
	result        *RtuOperationResult
	data          DataSet
	backgroundJob bool
	state         string
}

// RtuOperationResult
type RtuOperationResult struct {
	Name     string
	ID       string
	Result   interface{}
	Errors   []error
	Duration time.Duration
}

func (rtuo *RtuOperation) Start() {
	rtuo.started = time.Now()
	rtuo.state = "RUNNING"

	res, err := rtuo.operation(rtuo.client, rtuo.data)

	rtuo.result.Result = res
	rtuo.result.Duration = rtuo.Complete()
	rtuo.result.Name = rtuo.Name
	rtuo.result.ID = rtuo.ID

	if err != nil {
		rtuo.result.Errors = append(rtuo.result.Errors, err)
	}
}

func (rtuo *RtuOperation) Started() time.Time {
	return rtuo.started
}

func (rtuo *RtuOperation) Complete() time.Duration {
	rtuo.completed = time.Now()
	rtuo.state = "DONE"
	return rtuo.Duration()
}

func (rtuo *RtuOperation) Completed() time.Time {
	return rtuo.completed
}

func (rtuo *RtuOperation) Duration() time.Duration {

	if rtuo.completed.Sub(rtuo.started) <= 0 {
		return 0
	}

	return rtuo.completed.Sub(rtuo.started)
}

func (rtuo *RtuOperation) GetData() DataSet {
	return rtuo.data
}

func (rtuo *RtuOperation) GetResult() *RtuOperationResult {
	return rtuo.result
}

func (rtuo *RtuOperation) GetState() string {
	return rtuo.state
}

func (rtuo *RtuOperation) Cancel() {
	if rtuo.state == "QUEUED" {
		rtuo.state = "CANCELED"
	}
}

func (rtuo *RtuOperation) GetOutputJSON() []byte {
	get := rtuo.GetResult()

	newObject := &map[string]interface{}{
		"Completed": rtuo.Completed(),
		"Started":   rtuo.Started(),
		"Duration":  rtuo.Duration().String(),
		"Output":    get,
	}
	return ToJSON(newObject)
}

// Internal Functions

func newRtuOperation(n string, c *RmClient, d DataSet, o func(cli *RmClient, data DataSet) (interface{}, error)) *RtuOperation {

	e := []error{}
	r := &RtuOperationResult{
		Errors: e,
	}

	dbo := &RtuOperation{
		result:        r,
		created:       time.Now(),
		started:       time.Time{},
		completed:     time.Time{},
		Name:          n,
		client:        c,
		operation:     o,
		data:          d,
		ID:            NewGUID(),
		backgroundJob: false,
	}

	return dbo
}
