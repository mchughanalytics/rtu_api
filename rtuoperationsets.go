package rtuapi

import (
	"context"
	"encoding/json"
	"time"
)

// RTU OPERATION SETS

type RtuOperationSet struct {
	Operations      []*RtuOperation
	ctx             context.Context
	ID              string
	Name            string
	lookupOperation map[string]*RtuOperation
	post            func(dbos *RtuOperationSet) bool
	result          interface{}
}

type RtuOperationSetResult struct {
	Name       string
	ID         string
	Result     interface{}
	Errors     []error
	Duration   time.Duration
	Operations []*RtuOperationResult
}

func (rtuos *RtuOperationSet) Post() bool {
	return rtuos.post(rtuos)
}

func (rtuos *RtuOperationSet) GetOutputJSON() []byte {

	ops := make(map[string]interface{})

	for id, o := range rtuos.lookupOperation {

		newObject := &map[string]interface{}{
			"Completed": nil,
			"Started":   nil,
			"Duration":  nil,
			"Output":    nil,
		}

		_ = json.Unmarshal(o.GetOutputJSON(), newObject)

		ops[id] = newObject
	}

	output := &map[string]interface{}{
		"ID":         rtuos.ID,
		"Operations": ops,
	}

	return ToJSON(output)
}

// Internal Functions

func newRtuOperationSet(ctx context.Context) *RtuOperationSet {

	lookup := &map[string]*RtuOperation{}

	return &RtuOperationSet{
		ID:              NewGUID(),
		ctx:             ctx,
		Operations:      []*RtuOperation{},
		lookupOperation: *lookup,
	}
}
