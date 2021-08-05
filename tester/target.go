package tester

import "github.com/aws/aws-sdk-go-v2/service/ssm/types"

// TagNameTarget fulfills the tagParamBuilder interface for targets instances to be selected by the tag:Name
type TagNameTarget struct {
	tagNameValue string // the string tagNameValue that should be the
}

func (tnt TagNameTarget) buildTargetParameters() []types.Target {
	values := []string{tnt.tagNameValue}
	return []types.Target{{Key: stringPointer("tag:Name"), Values: values}}
}

// NewTagNameTarget returns and instance of TagNameTarget that can be supplied to the tester
// NewTagNameTarget expects a tagNameValue that should be the string value of the tag:Name of the instances to be targeted.
func NewTagNameTarget(tagNameValue string) TagNameTarget {
	return TagNameTarget{
		tagNameValue: tagNameValue,
	}
}
