package tester

import "testing"

func TestShellTestCase(t *testing.T) {
	cases := []struct {
		shellTestCase                  ShellTestCase
		expectedCommandParameterString string
	}{
		{
			NewTestCase("bash -c '</dev/tcp/google.com/443'", true, 3),
			`timeout 3 bash -c '</dev/tcp/google.com/443';if [ $? -eq 0 ];then $(exit 0);else $(exit 1);fi`,
		},
		{
			shellTestCase:                  NewTestCase("curl google.com", false, 5),
			expectedCommandParameterString: `timeout 5 curl google.com;if [ $? -eq 0 ];then $(exit 1);else $(exit 0);fi`,
		},
	}
	for _, v := range cases {
		// Todo Add more test cases
		t.Run("TestShellCommandBuildCommandParameters", func(t *testing.T) {
			if e, a := v.expectedCommandParameterString, v.shellTestCase.buildCommandString(); e != a {
				t.Errorf("Expected the command parameter string to be \n%s, but got \n%s", e, a)
			}
		})
	}
}
