package shell

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type Shell struct {
	storageType types.StorageType
	predefined  *map[string]string
}

func NewShell() action.Action {
	return &Shell{
		storageType: types.UNTAGGED,
		predefined:  nil,
	}
}

func (s *Shell) StorageType() types.StorageType {
	return s.storageType
}

func (s *Shell) Predefined() *map[string]string {
	return s.predefined
}

func (s *Shell) Name() string {
	return "shell"
}

func (s *Shell) Description() string {
	filepath := "shell.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (s *Shell) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	command := payload
	log.Printf("Executing command: %s", command)

	cmd := exec.CommandContext(context.Background(), "/bin/sh", "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Command error: %v", err)
		exitCode := cmd.ProcessState.ExitCode()
		return formatOutput(stdout.String(), stderr.String(), exitCode)
	}

	return formatOutput(stdout.String(), stderr.String(), 0)
}

func formatOutput(stdout, stderr string, exitCode int) string {
	result := stdout
	if stderr != "" {
		result += fmt.Sprintf("\nSTDERR: %s\n", stderr)
	}
	if exitCode != 0 {
		result += fmt.Sprintf("\nEXIT CODE: %d", exitCode)
	}
	return result
}

func (s *Shell) Timeout() *time.Duration {
	return nil
}

func (s *Shell) ExamplePayload() *string {
	p := "ls -la"
	return &p
}

func (s *Shell) ExampleAttributes() map[string]string {
	return nil
}

func (s *Shell) RequiredVariables() []*string {
	return nil
}

func (s *Shell) RequiresUserConfirmation() bool {
	return true
}

func (s *Shell) GetNamespace() types.NamespaceType {
	return types.SHELL
}
