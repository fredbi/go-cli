package gflag

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var (
	_ pflag.Value      = &SliceValue[string]{}
	_ pflag.SliceValue = &SliceValue[string]{}
	_ pflag.Value      = &SliceValue[time.Duration]{}
	_ pflag.SliceValue = &SliceValue[time.Duration]{}

	rmQuote = strings.NewReplacer(`"`, "", `'`, "", "`", "")
)

type (
	/*
		// FlaggableTypes is a type constraint that holds slices of types supported by pflag.
		FlaggableTypes interface {
			[]time.Duration |
				[]net.IP |
				[]net.IPNet |
				[]net.IPMask |
				// for extended types
				~[]struct{}
		}

		FlaggablePrimitives interface {
			~[]string |
				~[]bool |
				~[]int |
				~[]int8 |
				~[]int16 |
				~[]int32 |
				~[]int64 |
				~[]uint |
				// no ~[]uint8, which is the same as ~[]byte and interpreted as a value, not a slice
				~[]uint16 |
				~[]uint32 |
				~[]uint64 |
				~[]float32 |
				~[]float64
		}
	*/

	// SliceValue is a generic type that implements github.com/spf13/pflag.Value and SliceValue.
	SliceValue[T FlaggablePrimitives | FlaggableTypes] struct {
		Value   *[]T
		changed bool
	}
)

// NewFlagSliceValue constructs a generic flag compatible with github.com/spf13/pflag.SliceValue.
//
// Since the flag type is inferred from the underlying data type, some flexibility allowed by pflag is not
// always possible at this point.
//
// For example, when T = []string, NewFlagSliceValue adopts the semantics of the pflag.StringSlice (with comma-separated values),
// whereas pflag also supports a StringArray flag.
func NewFlagSliceValue[T FlaggablePrimitives | FlaggableTypes](addr *[]T, defaultValue []T) *SliceValue[T] {
	m := &SliceValue[T]{
		Value: addr,
	}
	*m.Value = defaultValue

	return m
}

func (m *SliceValue[T]) String() string {
	return writeAsSlice(m.GetSlice())
}

// Set knows how to config a string representation of the Value into a type T.
func (m *SliceValue[T]) Set(strValue string) error {
	slice, err := readAsCSV(rmQuote.Replace(strValue))
	if err != nil {
		return err
	}

	if !m.changed {
		// replace any default value
		err = m.Replace(slice)
		if err != nil {
			return err
		}

		m.changed = true

		return nil
	}

	// handle multiple occurences of the same flag with append semantics
	if err = m.append(slice...); err != nil {
		return err
	}

	return nil
}

func (m *SliceValue[T]) Type() string {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.Type()
	case *[]string:
		return "stringSlice"
	case *[]bool:
		return "boolSlice"
	case *[]int:
		return "intSlice"
	case *[]int8:
		return "int8Slice"
	case *[]int16:
		return "int16Slice"
	case *[]int32:
		return "int16Slice"
	case *[]int64:
		return "int64Slice"
	case *[]uint:
		return "uintSlice"
	case *[]uint16:
		return "uint16Slice"
	case *[]uint32:
		return "uint32Slice"
	case *[]uint64:
		return "uint64Slice"
	case *[]float32:
		return "float32Slice"
	case *[]float64:
		return "float64Slice"
	case *[]time.Duration:
		return "durationSlice"
	case *[]net.IP:
		return "ipSlice"
	case *[]net.IPNet:
		return "ipNetSlice"
	case *[]net.IPMask:
		return "ipMaskSlice"
	default:
		return fmt.Sprintf("%T", *m.Value)
	}
}

// Append a single element to the SliceValue, from its string representation.
func (m *SliceValue[T]) Append(strValue string) error {
	asAny := any(m.Value)
	if sliceValue, ok := asAny.(pflag.SliceValue); ok {
		return sliceValue.Append(strValue)
	}

	return m.append(strValue)
}

func (m *SliceValue[T]) Replace(strValues []string) error {
	asAny := any(m.Value)
	if sliceValue, ok := asAny.(pflag.SliceValue); ok {
		return sliceValue.Replace(strValues)
	}

	*m.Value = (*m.Value)[:0]

	return m.append(strValues...)
}

func (m *SliceValue[T]) append(strValues ...string) error {
	asAny := any(m.Value)

	switch v := asAny.(type) {
	case *[]string:
		*v = append(*v, strValues...)
		*m.Value = *cast[[]T](v)
	case *[]bool:
		slice, err := buildSliceFromParser(strValues, strconv.ParseBool)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int:
		slice, err := buildSliceFromParser(strValues, intParser[int](0))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int8:
		slice, err := buildSliceFromParser(strValues, intParser[int8](8))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int16:
		slice, err := buildSliceFromParser(strValues, intParser[int16](16))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int32:
		slice, err := buildSliceFromParser(strValues, intParser[int32](32))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int64:
		slice, err := buildSliceFromParser(strValues, intParser[int64](64))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint:
		slice, err := buildSliceFromParser(strValues, uintParser[uint](0))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint16:
		slice, err := buildSliceFromParser(strValues, uintParser[uint16](16))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint32:
		slice, err := buildSliceFromParser(strValues, uintParser[uint32](32))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint64:
		slice, err := buildSliceFromParser(strValues, uintParser[uint64](64))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]float32:
		slice, err := buildSliceFromParser(strValues, floatParser[float32](32))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]float64:
		slice, err := buildSliceFromParser(strValues, floatParser[float64](64))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]time.Duration:
		slice, err := buildSliceFromParser(strValues, time.ParseDuration)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]net.IP:
		slice, err := buildSliceFromParser(strValues, parseIP)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]net.IPNet:
		slice, err := buildSliceFromParser(strValues, parseIPNet)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]net.IPMask:
		slice, err := buildSliceFromParser(strValues, parseIPMask)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

	return nil
}

// GetSlice return a []string representation of the slice values.
func (m *SliceValue[T]) GetSlice() []string {
	asAny := any(m.Value)

	switch v := asAny.(type) {
	case pflag.SliceValue:
		return v.GetSlice()
	case *[]string:
		return *v
	case *[]bool:
		return buildSliceWithFormatter(*v, strconv.FormatBool)
	case *[]int:
		return buildSliceWithFormatter(*v, formatInt[int])
	case *[]int8:
		return buildSliceWithFormatter(*v, formatInt[int8])
	case *[]int16:
		return buildSliceWithFormatter(*v, formatInt[int16])
	case *[]int32:
		return buildSliceWithFormatter(*v, formatInt[int32])
	case *[]int64:
		return buildSliceWithFormatter(*v, formatInt[int64])
	case *[]uint:
		return buildSliceWithFormatter(*v, formatUint[uint])
	case *[]uint16:
		return buildSliceWithFormatter(*v, formatUint[uint16])
	case *[]uint32:
		return buildSliceWithFormatter(*v, formatUint[uint32])
	case *[]uint64:
		return buildSliceWithFormatter(*v, formatUint[uint64])
	case *[]float32:
		return buildSliceWithFormatter(*v, floatFormatter[float32](32))
	case *[]float64:
		return buildSliceWithFormatter(*v, floatFormatter[float64](64))
	case *[]time.Duration:
		return buildSliceWithFormatter(*v, formatStringer[time.Duration])
	case *[]net.IP:
		return buildSliceWithFormatter(*v, formatStringer[net.IP])
	case *[]net.IPNet:
		return buildSliceWithFormatter(*v, ipnetFormatter)
	case *[]net.IPMask:
		return buildSliceWithFormatter(*v, formatStringer[net.IPMask])
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
}

// MarshalFlag implements go-flags Marshaller interface
func (m *SliceValue[T]) MarshalFlag() (string, error) {
	return m.String(), nil
}

// UnmarshalFlag implements go-flags Unmarshaller interface
func (m *SliceValue[T]) UnmarshalFlag(value string) error {
	value = strings.TrimPrefix(value, `[`)
	value = strings.TrimSuffix(value, `]`)

	return m.Set(value)
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}

	csvReader := csv.NewReader(strings.NewReader(val))

	res, err := csvReader.Read()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return res, nil
}

func writeAsCSV(values []string) (string, error) {
	b := &bytes.Buffer{}
	csvWriter := csv.NewWriter(b)

	if err := csvWriter.Write(values); err != nil {
		return "", err
	}

	csvWriter.Flush()

	return strings.TrimSuffix(b.String(), "\n"), csvWriter.Error()
}

func writeAsSlice(slice []string) string {
	out, _ := writeAsCSV(slice)

	return fmt.Sprintf("[%s]", out)
}

func buildSliceWithFormatter[T any](slice []T, formatter func(T) string) []string {
	rep := make([]string, 0, len(slice))
	for _, val := range slice {
		rep = append(rep, formatter(val))
	}

	return rep
}

func buildSliceFromParser[T any](slice []string, parser func(string) (T, error)) ([]T, error) {
	rep := make([]T, 0, len(slice))
	for _, strVal := range slice {
		val, err := parser(strVal)
		if err != nil {
			return nil, err
		}

		rep = append(rep, val)
	}

	return rep, nil
}

func formatStringer[T fmt.Stringer](in T) string {
	return in.String()
}

func ipnetFormatter(in net.IPNet) string {
	return in.String()
}
