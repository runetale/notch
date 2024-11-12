package planning

import (
	"os"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type Clear struct {
}

func NewClear() action.Action {
	return &Clear{}
}

func (a *Clear) Name() string {
	return "clear_plan"
}

func (a *Clear) Description() string {
	filepath := "clear.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (a *Clear) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	return "plan clear"
}

func (a *Clear) Timeout() *time.Duration {
	return nil
}

func (a *Clear) ExamplePayload() *string {
	p := "complete the task"
	return &p
}

func (a *Clear) ExampleAttributes() map[string]string {
	return nil
}

func (a *Clear) RequiredVariables() []*string {
	return nil
}

func (a *Clear) RequiresUserConfirmation() bool {
	return true
}

func (a *Clear) GetNamespace() types.NamespaceType {
	return types.PLANNING
}

func (a *Clear) NamespaceDescription() string {
	filepath := "ns.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}
