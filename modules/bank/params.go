package bank

import (
	yaml "gopkg.in/yaml.v2"
)

// String implements the stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// String implements stringer insterface
func (se SendEnabled) String() string {
	out, _ := yaml.Marshal(se)
	return string(out)
}

// String returns a human readable string representation of a supplier.
func (supply Supply) String() string {
	bz, _ := yaml.Marshal(supply)
	return string(bz)
}
