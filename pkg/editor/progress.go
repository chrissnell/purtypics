package editor

import "sync"

// ProgressTracker tracks the progress of long-running operations.
type ProgressTracker struct {
	Progress int
	Status   string
	Error    string
	running  bool
	mutex    sync.RWMutex
}

// NewProgressTracker creates a new progress tracker.
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		Status: "idle",
	}
}

// TryStart atomically claims the tracker for a new operation. Returns false
// if an operation is already in progress — callers should reject the new
// request with 409 Conflict rather than corrupting shared state.
func (pt *ProgressTracker) TryStart() bool {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	if pt.running {
		return false
	}
	pt.running = true
	pt.Progress = 0
	pt.Status = "starting"
	pt.Error = ""
	return true
}

// Finish clears the running flag. Safe to call even if Start was never called.
func (pt *ProgressTracker) Finish() {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	pt.running = false
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
	pt.running = false
}

// Reset resets the tracker to idle state.
func (pt *ProgressTracker) Reset() {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	pt.Progress = 0
	pt.Status = "idle"
	pt.Error = ""
	pt.running = false
}

// Get returns the current progress state.
func (pt *ProgressTracker) Get() (progress int, status string, err string) {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()
	return pt.Progress, pt.Status, pt.Error
}