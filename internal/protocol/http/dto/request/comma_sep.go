package request

import "strings"

type CommaSepArray string

func (c CommaSepArray) Values() []string {
	str := string(c)
	if len(str) == 0 {
		return nil
	}

	return strings.Split(str, ",")
}
