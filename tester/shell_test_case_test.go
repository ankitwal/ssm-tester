package tester

import (
	"reflect"
	"testing"
)

func TestShellTestCase(t *testing.T) {
	cases := []struct {
		shellTestCase             ShellTestCase
		expectedCommandParameters map[string][]string
	}{
		{
			NewShellTestCase("timeout 3 bash -c '</dev/tcp/google.com/443'", true),
			map[string][]string{
				"commands":         []string{`$(timeout 3 bash -c '</dev/tcp/google.com/443');if [ $? -eq 0 ];then $(exit 0);else $(exit 1);fi`},
			},
		},
		{
			shellTestCase: NewShellTestCase("curl google.com", false),
			expectedCommandParameters: map[string][]string{
				"commands":         []string{`$(curl google.com);if [ $? -eq 0 ];then $(exit 1);else $(exit 0);fi`},
			},
		},
	}
	for _, v := range cases {
		// Todo Add more test cases
		t.Run("TestShellCommandBuildCommandParameters", func(t *testing.T) {
			if e, a := v.expectedCommandParameters, v.shellTestCase.buildCommandParameters(); !reflect.DeepEqual(e, a) {
				t.Errorf("Expected the command parameter string to be \n%v, but got \n%v", e, a)
			}
		})
	}
}
