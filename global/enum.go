package global

// CleanupType represents all the different types of cleanup we can perform
type CleanupType int

const (
	// Unconfirmed cleans up unconfirmed comments
	Unconfirmed CleanupType = 0
	// Deleted cleans up deleted comments
	Deleted CleanupType = 1
)
