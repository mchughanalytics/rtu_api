package rtu_api

import "log"

// RtuState manages current application state
type RtuState struct {
	RmClient   *RmClient
	UIProcess  int
	API        *RmApi
	state      chan *RtuInternalStateChange
	Operations chan *RtuOperationSet
	keepAlive  bool
}

type RtuInternalStateChange struct {
	NewState string
	Data     DataSet
}

func newRtuInternalState(state string) *RtuInternalStateChange {
	return &RtuInternalStateChange{
		NewState: state,
		Data:     nil,
	}
}

func NewRtuState() *RtuState {

	rtus := &RtuState{}

	s := make(chan *RtuInternalStateChange, 10)
	rtus.state = s

	o := make(chan *RtuOperationSet, 50)
	rtus.Operations = o

	return rtus
}

func (rtus *RtuState) Init(ui int, ka bool, key string) {

	api := NewRmAPI("/api", ":5000", key)

	go api.Init(rtus.Operations, rtus.state)

	rtus.API = api

	rtus.UIProcess = ui
	rtus.keepAlive = ka

	// TEMP
	// cli, err := NewRmClient("10.11.99.1", "root", "9Cdoe1dZQs")
	// if err != nil {
	// 	panic(err)
	// }
	// rtus.RmClient = cli

	go rtus.ListenForStateChange()

	for rtus.keepAlive == true {

	}
}

func (rtus *RtuState) ProcessOperation() {
	for opSet := range rtus.Operations {

		// for each operation set
		for _, op := range opSet.Operations {

			// for each operation within the set
			if op.backgroundJob {
				go op.Start()
			} else {
				op.Start()
			}

		}

		// after all operations
		opSet.Post()

	}

	// after channel closes
}

func (rtus *RtuState) ListenForStateChange() {
	waitingForShutdown := false

	for sc := range rtus.state {

		if sc.NewState == "stop" {
			log.Println("Received STOP")
			break
		}

		if sc.NewState == "process_then_stop" {
			log.Println("Received PROCESS THEN STOP")
			waitingForShutdown = true
		}

		if sc.NewState == "processing_complete" && waitingForShutdown {
			log.Println("Received PROCESSING COMPLETE - SHUTTING DOWN")
			break
		}

		if sc.NewState == "processing_complete" && !waitingForShutdown {
			log.Println("Received PROCESSING COMPLETE")
		}

		if sc.NewState == "device_configure" {
			rtus.configureDevice(sc.Data)
		}

	}

	rtus.keepAlive = false
}

func (rtus *RtuState) configureDevice(config DataSet) {
	c := &RmDeviceConfiguration{}
	_ = TransformStructs(config, c)

	cli, err := NewRmClient(c.Host, c.Username, c.Password)
	if err != nil {
		log.Println(err)
	}

	rtus.RmClient = cli
}
