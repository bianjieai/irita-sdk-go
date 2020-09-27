package auth

import (
	yaml "gopkg.in/yaml.v2"
)

// String implements the stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
