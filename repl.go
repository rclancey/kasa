package kasa

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

func getMethod(dev reflect.Value, name string) (reflect.Value, reflect.Method) {
	t := dev.Type()
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		method := t.Method(i)
		if method.IsExported() && method.Name == name {
			return dev.Method(i), method
		}
	}
	return reflect.Value{}, reflect.Method{}
}

func makeVal(t reflect.Type, s string) (reflect.Value, error) {
	typeName := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	switch typeName {
	case "time.Duration":
		f, err := strconv.ParseFloat(s, 64)
		return reflect.ValueOf(time.Duration(float64(time.Second)*f)), err
	case "time.Time":
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return reflect.ValueOf(time.Unix(int64(f), int64(math.Remainder(f, 1)*1e9))), nil
		}
		layouts := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05.999999999Z07:00",
			"2006-01-02 15:04:05.999999999 MST",
			"01/02/2006 15:04:05.999999999Z07:00",
			"01/02/2006 15:04:05.999999999 MST",
			"01/02/2006 3:04:05PM -0700",
			"01/02/2006 3:04:05PM MST",
		}
		for _, layout := range layouts {
			tm, err := time.Parse(layout, s)
			if err == nil {
				return reflect.ValueOf(tm), err
			}
		}
		return reflect.Value{}, fmt.Errorf("can't parse time value '%s'", s)
	}
	switch t.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		return reflect.ValueOf(b), err
	case reflect.Int:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(i), err
	case reflect.Int8:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(int8(i)), err
	case reflect.Int16:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(int16(i)), err
	case reflect.Int32:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(int32(i)), err
	case reflect.Int64:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(int64(i)), err
	case reflect.Uint:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(uint(i)), err
	case reflect.Uint8:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(uint8(i)), err
	case reflect.Uint16:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(uint16(i)), err
	case reflect.Uint32:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(uint32(i)), err
	case reflect.Uint64:
		i, err := strconv.Atoi(s)
		return reflect.ValueOf(uint64(i)), err
	case reflect.Float32:
		f, err := strconv.ParseFloat(s, 32)
		return reflect.ValueOf(float32(f)), err
	case reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		return reflect.ValueOf(f), err
	case reflect.String:
		return reflect.ValueOf(s), nil
	}
	return reflect.Value{}, fmt.Errorf("can't convert to %s", typeName)
}

func makeArgVals(method reflect.Method, args []string) ([]reflect.Value, error) {
	n := method.Type.NumIn() - 1
	argVals := make([]reflect.Value, n)
	for i := range argVals {
		t := method.Type.In(i+1)
		if method.Type.IsVariadic() && i == n - 1 {
			xvals := reflect.MakeSlice(t.Elem(), len(args) + 1 - n, len(args) + 1 - n)
			for j, arg := range args[i:] {
				val, err := makeVal(t.Elem(), arg)
				if err != nil {
					return argVals, err
				}
				xvals.Index(j).Set(val)
			}
			argVals[i] = xvals
		} else {
			val, err := makeVal(t, args[i])
			if err != nil {
				return argVals, err
			}
			argVals[i] = val
		}
	}
	return argVals, nil
}

func ExecDeviceCommand(dev SmartDevice, command string, args ...string) ([]any, error) {
	rdev := reflect.ValueOf(dev)
	fnc, method := getMethod(rdev, command)
	var zero reflect.Value
	if fnc == zero {
		return nil, fmt.Errorf("Unknown method '%s'", command)
	}
	argVals, err := makeArgVals(method, args)
	if err != nil {
		return nil, err
	}
	out := fnc.Call(argVals)
	resp := make([]any, len(out))
	for i, val := range out {
		resp[i] = val.Interface()
	}
	return resp, nil
}
