package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func convertSetToArray(set *schema.Set) []string {
	setList := set.List()
	retVal := []string{}
	for _, s := range setList {
		line := ""
		if s != nil {
			line = s.(string)
		}
		retVal = append(retVal, line)
	}
	return retVal
}
