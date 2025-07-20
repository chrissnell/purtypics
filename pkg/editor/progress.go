package editor

import "sync"

// ProgressTracker tracks the progress of long-running operations.
type ProgressTracker struct {
	Progress int
	Status   string
	Error    string
	mutex    sync.RWMutex
}

// NewProgressTracker creates a new progress tracker.
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		Status: "idle",
	}
}

// Update updates the progress and status.
func (pt *ProgressTracker) Update(progress int, status string) {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	pt.Progress = progress
	pt.Status = status
}

// SetError sets an error message and marks status as error.
func (pt *ProgressTracker) SetError(err string) {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	pt.Error = err
	pt.Status = "error"
}

// Reset resets the tracker to idle state.
func (pt *ProgressTracker) Reset() {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	pt.Progress = 0
	pt.Status = "idle"
	pt.Error = ""
}

// Get returns the current progress state.
func (pt *ProgressTracker) Get() (progress int, status string, err string) {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()
	return pt.Progress, pt.Status, pt.Error
}