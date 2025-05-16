package internal

import "strings"

type StringEnum struct {
	members []string
}

func (d *StringEnum) Members() []string {
	return d.members
}

func (d *StringEnum) IsMember(id string) bool {
	id = strings.ToLower(id)
	for _, member := range d.members {
		if member == id {
			return true
		}
	}
	return false
}

var DebtTypes = StringEnum{
	members: []string{"code", "documentation", "testing", "architecture", "infrastructure", "security"},
}

var DebtStatus = StringEnum{
	members: []string{"pending", "remediated", "in_progress"},
}
