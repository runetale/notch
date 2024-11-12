package tasklet

import (
	"os"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type Tasklet struct {
	name             string
	description      string
	workingDirectory string
	maxShownOutput   uint32
	args             map[string]string
	examplePayload   *string
	timeout          string
	tool             string
	storageType      types.StorageType
	predefined       *map[string]string
}

// TODO: implement complete and impossible
func NewTasklet() action.Action {
	return &Tasklet{
		name:        "tasklet",
		storageType: types.UNTAGGED,
		predefined:  nil,
	}
}

func (s *Tasklet) Name() string {
	return s.name
}

func (s *Tasklet) ExamplePayload() *string {
	p := "brief report on why the task is not possible"
	return &p
}

func (s *Tasklet) ExampleAttributes() map[string]string {
	return nil
}

func (s *Tasklet) Predefined() *map[string]string {
	return s.predefined
}

func (s *Tasklet) StorageType() types.StorageType {
	return s.storageType
}

func (s *Tasklet) Description() string {
	filepath := "shell.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (s *Tasklet) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	return "run"
}

func (s *Tasklet) Timeout() *time.Duration {
	return nil
}

func (s *Tasklet) RequiredVariables() []*string {
	return nil
}

func (s *Tasklet) RequiresUserConfirmation() bool {
	return true
}

func (s *Tasklet) GetNamespace() types.NamespaceType {
	return types.TASKLET
}

func (s *Tasklet) NamespaceDescription() string {
	filepath := "ns.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}
