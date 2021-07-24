package tester

import "fmt"

type ShellTestCase struct {
	command          string // the shell command to run as the test
	condition        bool   // to check for exit code 0 if true or not 0 if false
	executionTimeout int    // timeout in seconds for the the command to complete it
}

func (stc ShellTestCase) documentName() string {
	return "AWS-RunShellScript"
}

func (stc ShellTestCase) documentVersion() string {
	return "$LATEST"
}

func (stc ShellTestCase) buildCommandString() string {
	const stringTemplatePrefix = `timeout %d %s;if [ $? -eq 0 ];then $(exit %d);else $(exit %d);fi`
	exitCodeForSuccess := 0
	exitCodeForFailure := 1
	// if test condition is false, i.e check for failure, then swap the exit codes
	if !stc.condition {
		exitCodeForFailure, exitCodeForSuccess = exitCodeForSuccess, exitCodeForFailure
	}
	result := fmt.Sprintf(stringTemplatePrefix, stc.executionTimeout, stc.command, exitCodeForSuccess, exitCodeForFailure)
	fmt.Println(result)
	return result
}

func (stc ShellTestCase) buildCommandParameters() map[string][]string {
	const (
		commands         = "commands"
		executionTimeout = "executionTimeout"
	)
	parameters := map[string][]string{}
	parameters[commands] = []string{stc.buildCommandString()}
	parameters[executionTimeout] = []string{fmt.Sprintf("%d", stc.executionTimeout)}
	return parameters
}

// Todo add documentation
func NewTestCase(command string, condition bool, executionTimeout int) ShellTestCase {
	return ShellTestCase{
		command:          command,
		condition:        condition,
		executionTimeout: executionTimeout,
	}
}
