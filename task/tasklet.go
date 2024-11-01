package task

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Task interface {
	// これらの関数を表示
	// fn to_system_prompt(&self) -> Result<String>;
	// fn to_prompt(&self) -> Result<String>;
	// fn get_functions(&self) -> Vec<Namespace>;

	// fn get_timeout(&self) -> Option<Duration> {
	// 	None
	// }

	// fn get_rag_config(&self) -> Option<mini_rag::Configuration> {
	// 	None
	// }

	// fn max_history_visibility(&self) -> u16 {
	// 	50
	// }

	// fn guidance(&self) -> Result<Vec<String>> {
	// 	self.base_guidance()
	// }

	// fn namespaces(&self) -> Option<Vec<String>> {
	// 	None
	// }

	//	fn base_guidance(&self) -> Result<Vec<String>> {
	//		// basic rules to extend
	//		Ok(include_str!("basic_guidance.prompt")
	//			.split('\n')
	//			.map(|l| l.trim().to_string())
	//			.filter(|l| !l.is_empty())
	//			.collect())
	//	}
}

type Tasklet struct {
	Name         string
	Folder       string
	Using        []*string  `yaml:"using"`
	SystemPrompt string     `yaml:"system_prompt"`
	Prompt       *string    `yaml:"prompt"`
	Guidance     []string   `yaml:"guidance"`
	Functions    []Function `yaml:"functions"`
}

type Function struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Actions     []Action `yaml:"actions"`
}

type Action struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	Tool           string `yaml:"tool"`
	MaxShownOutput int    `yaml:"max_shown_output"`
	ExamplePayload string `yaml:"example_payload,omitempty"`
}

func GetFromPath(path string) (*Tasklet, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return getFromDir(path)
	}
	return getFromYamlFile(path)
}

func getFromDir(path string) (*Tasklet, error) {
	filePath := filepath.Join(path, "task.yaml")
	_, err := os.Stat(filePath)
	if err == nil {
		return getFromYamlFile(filePath)
	}
	if os.IsNotExist(err) {
		return nil, err
	}

	return nil, err
}

func getFromYamlFile(filePath string) (*Tasklet, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("can't read file %v", err)
		return nil, err
	}

	var tasklet Tasklet
	err = yaml.Unmarshal(data, &tasklet)
	if err != nil {
		log.Fatalf("parsing yaml error %v", err)
		return nil, err
	}

	dir := filepath.Dir(filePath)
	dirName := filepath.Base(dir)
	tasklet.Name = dirName
	tasklet.Folder = filePath

	return &tasklet, nil

}

// userPromptが無ければ入力を受け付ける
// Promptが設定されていない場合はuserからのpromptを設定する
func (t *Tasklet) Setup(userPrompt *string) error {
	if userPrompt == nil {
		input := getUserInput("enter task > ")
		t.Prompt = &input
		return nil
	}

	if t.Prompt == nil {
		t.Prompt = userPrompt
		return nil
	}

	return errors.New("Setup failed")
}

func (t *Tasklet) GetUsing() []*string {
	return t.Using
}

// user defined yaml tasks
func (t *Tasklet) GetFunctions() []Function {
	fs := t.Functions
	fmt.Println("Called GetFunctions")
	fmt.Println(fs)
	return t.Functions
}

func getUserInput(prompt string) string {
	fmt.Print("\n" + prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	fmt.Println()
	return strings.TrimSpace(input)
}
