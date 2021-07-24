package tester

import "github.com/aws/aws-sdk-go-v2/service/ssm/types"

// Interface type to abstract concrete type from tester
type targetParamBuilder interface {
	buildTargetParameters() []types.Target
}

// Target tupe for targets with tag name
type tagNameTarget struct {
	tagNameValue string // the string tagNameValue that should be the
}

func (tnt tagNameTarget) buildTargetParameters() []types.Target {
	values := []string{tnt.tagNameValue}
	return []types.Target{{Key: stringPointer("tag:Name"), Values: values}}
}

// Todo add documentation
func NewTagNameTarget(tagNameValue string) tagNameTarget {
	return tagNameTarget{
		tagNameValue: tagNameValue,
	}
}
// Todo add target type to identify ec2 instances with multiple tags / instance IDs
