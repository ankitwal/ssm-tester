package tester

import "fmt"

type ShellTestCase struct {
	command          string // the shell command to run as the test
	condition        bool   // to check for exit code 0 if true or not 0 if false
	executionTimeout int    // timeout in seconds for the the command to complete it
}

// Todo add documentation
func NewTestCase(command string, condition bool, executionTimeout int) *TestCase {
	return &ShellTestCase{
		command:          command,
		condition:        condition,
		executionTimeout: executionTimeout,
	}
}

func (t *ShellTestCase) buildCommandParameters() string {
	const stringTemplatePrefix = `if $(timeout %d '%s');then exit %d;else exit %d;fi`
	exitCode := 1
	if t.condition {
		exitCode = 0
	}
	result := fmt.Sprintf(stringTemplatePrefix, t.executionTimeout, t.command, exitCode, exitCode)
	fmt.Println(result)
	return result
}
