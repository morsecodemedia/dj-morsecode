package mpv

import "testing"

func TestNewProcess(t *testing.T) {

	process := NewProcess(
		"/tmp/dj-morsecode-test.sock",
	)

	if process == nil {
		t.Fatal("expected process")
	}

	if process.socketPath != "/tmp/dj-morsecode-test.sock" {
		t.Errorf(
			"expected socket path %q, got %q",
			"/tmp/dj-morsecode-test.sock",
			process.socketPath,
		)
	}

	if process.cmd != nil {
		t.Fatal("expected process not to be started")
	}

}

func TestProcessStopBeforeStart(t *testing.T) {

	process := NewProcess(
		"/tmp/dj-morsecode-test.sock",
	)

	if err := process.Stop(); err != nil {
		t.Fatalf(
			"expected stop before start to succeed: %v",
			err,
		)
	}

}
