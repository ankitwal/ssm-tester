package tester

import "fmt"

// ShellTestCase configuration for a shell script command and a condition
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

// NewShellTestCase is a constructor for ShellTestCase type.
// It expects command and condition arguments.
// command is a string representation of the shell test command that should be run in a test. eg. "echo foo" or "timeout 2 bash -c '</dev/tcp/google.com/443'" to
// test tcp connectivity to google.com:443 in 2 seconds.
// condition is bool that represents if the test should check of success(i.e exit code is 0) or failure(i.e exit code is not 0).
//
// Important Note - validity of the command depends on the OS and Binaries installed on the target instances.
// It is up to the user to ensure that that is the command uses a binary, it is installed on the target instance.
// Eg. if the command is 'aws s3 ls' then the user should ensure that the aws cli is installed on the target instance of the test would always fail.
func NewShellTestCase(command string, condition bool) ShellTestCase {
	return ShellTestCase{
		command:   command,
		condition: condition,
	}
}
