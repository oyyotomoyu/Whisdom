// Package logs implements Whisdom's audit/system log: one JSON Lines file
// per day, written through a small typed API so callers can never choose
// their own filename or an invalid status.
package logs

import "fmt"

// Status levels accepted by the log service. Any other value is rejected.
const (
	StatusInfo    = "info"
	StatusWarning = "warning"
	StatusError   = "error"
)

const timestampLayout = "2006-01-02 15:04:05 -07:00"

// Record is a single log line.
type Record struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	IP        string `json:"ip"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
}

func validateStatus(status string) error {
	switch status {
	case StatusInfo, StatusWarning, StatusError:
		return nil
	default:
		return fmt.Errorf("invalid log status %q", status)
	}
}
