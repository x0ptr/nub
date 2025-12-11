package main

import (
	"os"
	"os/exec"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

type execCmd struct {
	name string
	args []string
}

func (e *execCmd) run() error {
	return runCommand(e.name, e.args...)
}

func (e *execCmd) runWait() error {
	cmd := exec.Command(e.name, e.args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
