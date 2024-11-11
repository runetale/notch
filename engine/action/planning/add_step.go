package planning

import (
	"os"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type AddStep struct {
}

func NewAddStep() action.Action {
	return &AddStep{}
}

func (a *AddStep) Name() string {
	return "add_plan_step"
}

func (a *AddStep) Description() string {
	filepath := "add.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (a *AddStep) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	storage.AddCompletion(payload)
	return "step added to the plan"
}

func (a *AddStep) Timeout() *time.Duration {
	return nil
}

func (a *AddStep) ExamplePayload() *string {
	p := "complete the task"
	return &p
}

func (a *AddStep) ExampleAttributes() map[string]string {
	return nil
}

func (a *AddStep) RequiredVariables() []*string {
	return nil
}

func (a *AddStep) RequiresUserConfirmation() bool {
	return true
}

func (a *AddStep) GetNamespace() types.NamespaceType {
	return types.PLANNING
}

func (a *AddStep) NamespaceDescription() string {
	filepath := "ns.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}
