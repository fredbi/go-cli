package gflag

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/spf13/pflag"
	"golang.org/x/exp/constraints"
)

var (
	// type guards: Value implements pflag.Value
	_ pflag.Value = &Value[string]{}
	_ pflag.Value = &Value[time.Duration]{}
)

type (
	// FlaggableTypes is a type constraint that holds all types supported by pflag, besides primitive types.
	FlaggableTypes interface {
		time.Duration |
			net.IP |
			net.IPNet |
			net.IPMask |
			// for extended types
			~struct{}
	}

	// FlaggablePrimitives is a type constraint that holds all primitive types supported by pflag.
	//
	// Exception: complex types are not supported.
	FlaggablePrimitives interface {
		constraints.Integer |
			constraints.Float |
			~string |
			~bool |
			~[]byte // aka: ~[]uint8 |
	}

	// Value is a generic type that implements github.com/spf13/pflag.Value.
	Value[T FlaggablePrimitives | FlaggableTypes] struct {
		Value       *T
		NoOptDefVal string
	}
)

// NewFlagValue constructs a generic flag compatible with github.com/spf13/pflag.Value.
//
// Since the flag type is inferred from the underlying data type, some flexibility allowed by pflag is not
// always possible at this point.
//
// For example, when T = []byte, NewFlagValue adopts the semantics of the pflag.BytesHex flag, whereas pflag aslo supports
// a BytesBase64 flag.
//
// Similarly, when T = int, we adopt the semantics of pflag.Int and not pflag.Count.
func NewFlagValue[T FlaggablePrimitives | FlaggableTypes](addr *T, defaultValue T) *Value[T] {
	if addr == nil {
		panic("NewFlagValue must take a valid pointer to T")
	}

	m := &Value[T]{
		Value: addr,
	}
	*m.Value = defaultValue

	// bool flag imply true when set without arg
	v := any(m.Value)
	if _, isBool := v.(*bool); isBool {
		m.NoOptDefVal = "true"
	}

	return m
}

// GetValue returns the underlying value of the flag.
func (m Value[T]) GetValue() T {
	return *m.Value
}

// String knows how to yield a string representation of type T.
func (m *Value[T]) String() string {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.String()
	case *string:
		return *v
	case *int:
		return formatInt(*v)
	case *int8:
		return formatInt(*v)
	case *int16:
		return formatInt(*v)
	case *int32:
		return formatInt(*v)
	case *int64:
		return formatInt(*v)
	case *uint:
		return formatUint(*v)
	case *uint8:
		return formatUint(*v)
	case *uint16:
		return formatUint(*v)
	case *uint32:
		return formatUint(*v)
	case *uint64:
		return formatUint(*v)
	case *float32:
		return floatFormatter[float32](32)(*v)
	case *float64:
		return floatFormatter[float64](64)(*v)
	case *bool:
		return strconv.FormatBool(*v)
	case *[]byte:
		return fmt.Sprintf("%X", *v)
	case *time.Duration:
		return v.String()
	case *net.IP:
		return v.String()
	case *net.IPNet:
		return v.String()
	case *net.IPMask:
		return v.String()
	case fmt.Stringer:
		return v.String()
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
}

// Set knows how to config a string representation of the Value into a type T.
func (m *Value[T]) Set(strValue string) error {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.Set(strValue)
	case *string:
		val := strValue
		*m.Value = *cast[T](&val)
	case *int:
		val, err := intParser[int](0)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int8:
		val, err := intParser[int8](8)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int16:
		val, err := intParser[int16](16)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int32:
		val, err := intParser[int32](32)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int64:
		val, err := intParser[int64](64)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint:
		val, err := uintParser[uint](0)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint8:
		val, err := uintParser[uint8](8)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint16:
		val, err := uintParser[uint16](16)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint32:
		val, err := uintParser[uint32](32)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint64:
		val, err := uintParser[uint64](64)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *float32:
		val, err := floatParser[float32](32)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *float64:
		val, err := floatParser[float64](64)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *bool:
		val, err := strconv.ParseBool(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
		m.NoOptDefVal = "true"
	case *[]byte:
		val, err := hex.DecodeString(strings.TrimSpace(strValue))
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *time.Duration:
		val, err := time.ParseDuration(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *net.IP:
		val, err := parseIP(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *net.IPMask:
		val, err := parseIPMask(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *net.IPNet:
		val, err := parseIPNet(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

	return nil
}

func (m *Value[T]) Type() string {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.Type()
	case *time.Duration:
		return "duration"
	case *net.IP:
		return "ip"
	case *net.IPMask:
		return "ipMask"
	case *net.IPNet:
		return "ipNet"
	case *[]byte:
		return "byteHex"
	default:
		return fmt.Sprintf("%T", *m.Value)
	}
}

// MarshalFlag implements go-flags Marshaller interface
func (m *Value[T]) MarshalFlag() (string, error) {
	return m.String(), nil
}

// UnmarshalFlag implements go-flags Unmarshaller interface
func (m *Value[T]) UnmarshalFlag(value string) error {
	return m.Set(value)
}

func cast[U any, T any](v *T) *U {
	return (*U)(unsafe.Pointer(v))
}

func formatInt[T constraints.Signed](in T) string {
	return strconv.FormatInt(int64(in), 10)
}

func formatUint[T constraints.Unsigned](in T) string {
	return strconv.FormatUint(uint64(in), 10)
}

func floatFormatter[T constraints.Float](bits int) func(T) string {
	return func(in T) string {
		return strconv.FormatFloat(float64(in), 'g', -1, bits)
	}
}

func intParser[T constraints.Signed](bits int) func(in string) (T, error) {
	return func(in string) (T, error) {
		val, err := strconv.ParseInt(in, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

func uintParser[T constraints.Unsigned](bits int) func(string) (T, error) {
	return func(in string) (T, error) {
		val, err := strconv.ParseUint(in, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

func floatParser[T constraints.Float](bits int) func(string) (T, error) {
	return func(in string) (T, error) {
		val, err := strconv.ParseFloat(in, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

func parseIP(in string) (net.IP, error) {
	val := net.ParseIP(strings.TrimSpace(in))
	if val == nil {
		return nil, fmt.Errorf("failed to parse IP: %q", in)
	}

	return val, nil
}

func parseIPMask(in string) (net.IPMask, error) {
	val := pflag.ParseIPv4Mask(strings.TrimSpace(in))
	if val == nil {
		return nil, fmt.Errorf("failed to parse IP mask: %q", in)
	}

	return val, nil
}

func parseIPNet(in string) (net.IPNet, error) {
	_, val, err := net.ParseCIDR(strings.TrimSpace(in))
	if val == nil {
		return net.IPNet{}, fmt.Errorf("failed to parse CIDR: %q", in)
	}

	return *val, err
}
