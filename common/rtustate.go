package common

// RtuState manages current application state
type RtuState struct {
	RmClient  *RmClient
	UIProcess int
}
