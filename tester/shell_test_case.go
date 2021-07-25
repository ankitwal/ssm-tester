package tester

import "fmt"

// Todo Add documentation
type ShellTestCase struct {
	command   string // the shell command to run as the test
	condition bool   // to check for exit code 0 if true or not 0 if false
}

func (stc ShellTestCase) documentName() string {
	return "AWS-RunShellScript"
}

func (stc ShellTestCase) documentVersion() string {
	return "$LATEST"
}

func (stc ShellTestCase) buildCommandString() string {
	const stringTemplatePrefix = `$(%s);if [ $? -eq 0 ];then $(exit %d);else $(exit %d);fi`
	exitCodeForSuccess := 0
	exitCodeForFailure := 1
	// if test condition is false, i.e check for failure, then swap the exit codes
	if !stc.condition {
		exitCodeForFailure, exitCodeForSuccess = exitCodeForSuccess, exitCodeForFailure
	}
	result := fmt.Sprintf(stringTemplatePrefix, stc.command, exitCodeForSuccess, exitCodeForFailure)
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
	// the timeout rules for ssm look rather convoluted, this execution time out probably cannot be relied on
	// hence set to a higher number and will rely on consumer supply a timeout in the command
	// parameters[executionTimeout] = []string{"3600"}
	return parameters
}

// Todo add documentation
func NewShellTestCase(command string, condition bool) ShellTestCase {
	return ShellTestCase{
		command:   command,
		condition: condition,
	}
}
