package sinli

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

type sinliValue struct {
	length int
	value  reflect.Value
}

// Marshal converts a struct into a "sinli" formatted string.
func Marshal(v interface{}) ([]byte, error) {
	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		// Loop through the slice and marshal each element
		var output string
		for i := 0; i < value.Len(); i++ {
			marshaled, err := Marshal(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			output += string(marshaled)
		}
		return encode(output)
	}

	// If the value is a pointer, dereference it
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Ensure that the provided value is a struct
	if value.Kind() != reflect.Struct {
		return nil, errors.New("sinli: value must be a struct")
	}

	values := map[int]sinliValue{}
	var keys []int
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tag := value.Type().Field(i).Tag.Get("sinli")
		// Split the tag into its parts
		parts := strings.Split(tag, ",")

		var order, length int
		var fixed string
		for _, part := range parts {
			// Split the part into its key and value
			kv := strings.Split(part, "=")
			if len(kv) != 2 {
				return nil, errors.New("sinli: invalid tag format")
			}
			k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			switch k {
			case "order":
				candidate, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("sinli: invalid order '%s'", v)
				}
				order = candidate
			case "length":
				candidate, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("sinli: invalid length '%s'", v)
				}
				length = candidate
			case "fixed":
				fixed = v
			default:
				return nil, fmt.Errorf("sinli: invalid tag key '%s'", k)
			}
		}
		if order == 0 {
			return nil, errors.New("sinli: order must be specified")
		}
		if length == 0 && !isArray(field) && !isSinli(field) {
			return nil, fmt.Errorf("sinli: length must be specified for %v", field)
		}
		if _, ok := values[order]; ok {
			return nil, fmt.Errorf("sinli: duplicate order '%d'", order)
		}
		if fixed != "" {
			field = reflect.ValueOf(fixed)
		}
		values[order] = sinliValue{
			length: length,
			value:  field,
		}
		keys = append(keys, order)
	}
	var output string
	sort.Ints(keys)
	for _, key := range keys {
		kind := values[key].value.Kind()
		value := values[key].value
		length := values[key].length
		switch {
		case kind == reflect.Slice, kind == reflect.Array, isSinli(value):
			marshaled, err := Marshal(value.Interface())
			if err != nil {
				return nil, err
			}
			output += string(marshaled)
		default:
			output += toString(value, length)
		}
	}
	// Append a newline to the output if it doesn't already have one
	if !strings.HasSuffix(output, "\r\n") {
		output += "\r\n"
	}
	return encode(output)
}

func encode(v string) ([]byte, error) {
	// 850 OEM – Multilingual Latin I
	encoder := charmap.CodePage850.NewEncoder()
	latin, err := encoder.String(v)
	if err != nil {
		return nil, fmt.Errorf("sinli: couldn't encode text: %w", err)
	}
	return []byte(latin), nil
}

func isArray(v reflect.Value) bool {
	return v.Kind() == reflect.Slice || v.Kind() == reflect.Array
}

func isSinli(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("sinli")
		if tag != "" {
			return true
		}
	}
	return false
}

func toString(v reflect.Value, length int) string {
	// If the value is a pointer and it's nil, return an empty string
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return strings.Repeat(" ", length)
	}
	// Obtain the underlying value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "S"
		}
		return "N"
	// Check if the value is time.Time
	case reflect.Struct:
		if v.Type().String() == "time.Time" {
			return v.Interface().(time.Time).Format("20060102")
		}
	case reflect.Float32, reflect.Float64:
		f := roundFloat(v.Float(), 2)
		s := fmt.Sprintf("%0"+strconv.Itoa(length+1)+".2f", f)
		return strings.ReplaceAll(s, ".", "")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%0"+strconv.Itoa(length)+"d", v.Int())
	case reflect.String:
		// Convert string to
	}
	s := fmt.Sprintf("%v", v)
	return fmt.Sprintf("%-"+strconv.Itoa(length)+"s", s)
}

func roundFloat(number float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	rounded := math.Round(number*shift) / shift
	return rounded
}
