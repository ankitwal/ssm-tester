package tester

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"reflect"
	"testing"
)

func TestTagNameTarget(t *testing.T) {
	var cases = []struct {
		target   tagNameTarget
		expected []types.Target
	}{
		{
			target: NewTagNameTarget("dummyTagName"),
			expected: []types.Target{
				{
					Key:    stringPointer("tag:Name"),
					Values: []string{"dummyTagName"},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run("Ensure TagNameTarget Constructs SnedCommand Target Inputs Correctly", func(t *testing.T) {
			if e, a := c.expected, c.target.buildTargetParameters(); !reflect.DeepEqual(e, a) {
				t.Errorf("Expected %v, but got %v", e, a)
			}
		})
	}
}
