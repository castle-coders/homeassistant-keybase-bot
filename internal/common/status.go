package status

const (
	// Ready to indicate client is ready to recieve messages
	Ready = 0
	// Error to indicate that the client has failed
	Error = -1
	// Shutdown no longer processing messages
	Shutdown = 1
)
