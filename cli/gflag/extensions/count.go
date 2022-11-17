package extensions

import (
	"strconv"

	"github.com/spf13/pflag"
)

var _ pflag.Value = new(CountValue)

// CountValue is used to pass increment/decrement counts as flags.
type CountValue int

func NewCountValue(defaultValue int, p *int) *CountValue {
	*p = defaultValue
	return (*CountValue)(p)
}

func (i *CountValue) Set(s string) error {
	// "+1" means that no specific value was passed, so increment
	if s == "+1" {
		*i++

		return nil
	}

	v, err := strconv.ParseInt(s, 0, 0)
	*i = CountValue(v)

	return err
}

func (i *CountValue) Type() string {
	return "count"
}

func (i *CountValue) String() string {
	return strconv.Itoa(int(*i))
}

// CountVar like CountVar only the flag is placed on the CommandLine instead of a given flag set
func CountVar(p *int, name string, defaultValue int, usage string) {
	CountVarP(p, name, "", defaultValue, usage)
}

// CountVarP is like CountVar only take a shorthand for the flag name.
func CountVarP(p *int, name, shorthand string, defaultValue int, usage string) {
	v := NewCountValue(defaultValue, p)
	_ = pflag.CommandLine.VarPF(v, name, shorthand, usage)
}

// Count defines a count flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// A count flag will add 1 to its value evey time it is found on the command line
func Count(name string, defaultValue int, usage string) *int {
	return CountP(name, "", defaultValue, usage)
}

// CountP is like Count only takes a shorthand for the flag name.
func CountP(name, shorthand string, defaultValue int, usage string) *int {
	p := new(int)
	CountVarP(p, name, shorthand, defaultValue, usage)

	return p
}
