// xml pareser for state
package serializer

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/state"
	"github.com/runetale/notch/llm"
)

//go:embed actions.prompt
var actionPrompt string

//go:embed system.prompt
var systemPrompt string

type System struct {
	SystemPrompt     string `xml:"system_prompt"`
	Storages         string `xml:"storages"`
	Iterations       string `xml:"iterations"`
	AvailableActions string `xml:"available_actions"`
	Guidance         string `xml:"guidance"`
}

func DisplaySystemPrompt(state *state.State) (string, error) {
	// input data to template
	tmpl, err := template.New("prompt").ParseFiles(systemPrompt)
	if err != nil {
		return "", err
	}

	// system prompt
	task := state.GetTask()
	sysprompt := task.GetSystemPrompt()

	// storages
	storages := state.GetStorages()
	sortedStorageKeys := make([]string, 0, len(storages))
	for key := range storages {
		sortedStorageKeys = append(sortedStorageKeys, key)
	}
	sort.Strings(sortedStorageKeys)
	storage := strings.Join(sortedStorageKeys, "\n\n")

	// guidance
	var formattedGuidance []string
	for _, s := range task.Guidance {
		formattedGuidance = append(formattedGuidance, fmt.Sprintf("- %s", s))
	}
	guidance := strings.Join(formattedGuidance, "\n")

	// available actions
	actions, err := actionsForState(state)
	if err != nil {
		return "", err
	}
	availableActions := actionPrompt + "\n" + actions

	// iterations
	iterations := ""
	if state.GetMaxIteration() > 0 {
		iterations = fmt.Sprintf("You are currently at step %d of a maximum of %d", state.GetCurrentStep()+1, state.GetMaxIteration())
	}

	data := System{
		SystemPrompt:     sysprompt,
		Storages:         storage,
		Iterations:       iterations,
		AvailableActions: availableActions,
		Guidance:         guidance,
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}

	// 結果を出力
	return output.String(), nil
}

func serializeAction(ac action.Action) string {
	var builder strings.Builder

	// create xml tag
	builder.WriteString(fmt.Sprintf("<%s", ac.Name()))

	// if existing attributes by aciton
	for name, exampleValue := range ac.ExampleAttributes() {
		builder.WriteString(fmt.Sprintf(` %s="%s"`, name, exampleValue))
	}

	// if existing payload by aciton
	if payload := ac.ExamplePayload(); payload != nil {
		builder.WriteString(fmt.Sprintf(">%s</%s>", *payload, ac.Name()))
	} else {
		builder.WriteString("/>")
	}

	return builder.String()
}

func actionsForState(state *state.State) (string, error) {
	var builder strings.Builder

	for _, group := range state.GetNamespaces() {
		builder.WriteString(fmt.Sprintf("## %s\n\n", group.GetName()))

		if group.GetDescription() != "" {
			builder.WriteString(fmt.Sprintf("%s\n\n", group.GetDescription()))
		}

		for _, action := range group.GetActions() {
			builder.WriteString(fmt.Sprintf("%s %s\n\n",
				action.Description(),
				serializeAction(action),
			))
		}
	}

	return builder.String(), nil
}

func invocation(inv *llm.Invocation) string {
	var xml strings.Builder
	xml.WriteString(fmt.Sprintf("<%s", inv.Action))

	if inv.Attributes != nil {
		for key, value := range inv.Attributes {
			xml.WriteString(fmt.Sprintf(" %s=\"%s\"", key, value))
		}
	}

	payload := ""
	if inv.Payload != nil {
		payload = *inv.Payload
	}
	xml.WriteString(fmt.Sprintf(">%s</%s>", payload, inv.Action))

	return xml.String()
}

func SerializeInvocation(inv *llm.Invocation) *string {
	invocation := invocation(inv)
	return &invocation
}

func SerializeAction(ac *action.Action) string {
	return ""
}

func SerializeStorage(ac *action.Action) string {
	return ""
}

func TryParse(raw string) []*llm.Invocation {
	ptr := raw
	var parsedInvocations []*llm.Invocation
	uniqueMap := make(map[string]bool)

	for {
		openIdx := strings.Index(ptr, "<")
		if openIdx == -1 {
			break
		}

		ptr = ptr[openIdx:]

		parsedBlock := tryParseBlock(ptr)
		if parsedBlock.Processed == 0 {
			break
		}

		for _, inv := range parsedBlock.Invocations {
			uniqueKey := fmt.Sprintf("%s-%v-%v", inv.Action, inv.Attributes, inv.Payload)
			if !uniqueMap[uniqueKey] {
				uniqueMap[uniqueKey] = true
				parsedInvocations = append(parsedInvocations, inv)
			}
		}

		ptr = ptr[parsedBlock.Processed:]
	}

	return parsedInvocations
}
