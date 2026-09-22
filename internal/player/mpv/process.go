package mpv

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const (
	processReadyTimeout = 5 * time.Second
	processReadyPoll    = 50 * time.Millisecond
)

type Process struct {
	cmd        *exec.Cmd
	socketPath string
}

func NewProcess(
	socketPath string,
) *Process {

	return &Process{
		socketPath: socketPath,
	}

}

func (p *Process) Start() error {

	if p.cmd != nil {
		return fmt.Errorf(
			"mpv process already started",
		)
	}

	cmd := exec.Command(
		"mpv",
		"--idle=yes",
		"--no-terminal",
		"--input-ipc-server="+p.socketPath,
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf(
			"start mpv: %w",
			err,
		)
	}

	p.cmd = cmd

	return nil

}

func (p *Process) WaitReady() error {

	deadline := time.Now().Add(
		processReadyTimeout,
	)

	for time.Now().Before(deadline) {

		conn, err := net.DialTimeout(
			"unix",
			p.socketPath,
			processReadyPoll,
		)
		if err == nil {

			_ = conn.Close()

			return nil

		}

		time.Sleep(
			processReadyPoll,
		)

	}

	return fmt.Errorf(
		"mpv IPC socket not ready after %s",
		processReadyTimeout,
	)

}

func (p *Process) Stop() error {

	if p.cmd == nil {
		return nil
	}

	process := p.cmd.Process

	if process == nil {
		p.cmd = nil
		return p.removeSocket()
	}

	signalErr := process.Signal(
		os.Interrupt,
	)

	waitErr := p.cmd.Wait()

	p.cmd = nil

	socketErr := p.removeSocket()

	if signalErr != nil {
		return fmt.Errorf(
			"stop mpv: %w",
			signalErr,
		)
	}

	var exitErr *exec.ExitError

	if waitErr != nil &&
		!errors.As(waitErr, &exitErr) {

		return fmt.Errorf(
			"wait for mpv: %w",
			waitErr,
		)

	}

	if socketErr != nil {
		return socketErr
	}

	return nil

}

func (p *Process) removeSocket() error {

	err := os.Remove(
		p.socketPath,
	)

	if err == nil ||
		os.IsNotExist(err) {

		return nil
	}

	return fmt.Errorf(
		"remove mpv socket: %w",
		err,
	)

}
