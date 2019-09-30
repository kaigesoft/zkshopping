package tools

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func String(val interface{}) string {
	if val == nil {
		return ""
	}
	re_value := reflect.ValueOf(val)
	for re_value.Kind() == reflect.Ptr {
		re_value = re_value.Elem()
		if !re_value.IsValid() {
			return ""
		}
		val = re_value.Interface()
		if val == nil {
			return ""
		}
		re_value = reflect.ValueOf(val)
	}
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}

func Int64(val interface{}) (int64, bool) {
	if val == nil {
		return 0, false
	}
	re_value := reflect.ValueOf(val)
	for re_value.Kind() == reflect.Ptr {
		re_value = re_value.Elem()
		if !re_value.IsValid() {
			return 0, false
		}
		val = re_value.Interface()
		if val == nil {
			return 0, false
		}
		re_value = reflect.ValueOf(val)
	}
	if val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case bool:
		if v {
			return 1, true
		} else {
			return 0, true
		}
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	case []byte:
		return Int64(string(v))
	case string:
		v = strings.SplitN(v, ".", 2)[0]
		t, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return t, true
		} else {
			return 0, false
		}
	default:
		return 0, false
	}
}

func Uint64(val interface{}) (uint64, bool) {
	if val == nil {
		return 0, false
	}
	re_value := reflect.ValueOf(val)
	for re_value.Kind() == reflect.Ptr {
		re_value = re_value.Elem()
		if !re_value.IsValid() {
			return 0, false
		}
		val = re_value.Interface()
		if val == nil {
			return 0, false
		}
		re_value = reflect.ValueOf(val)
	}
	if val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case bool:
		if v {
			return 1, true
		} else {
			return 0, true
		}
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return uint64(v), true
	case int8:
		return uint64(v), true
	case int16:
		return uint64(v), true
	case int:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case float32:
		return uint64(v), true
	case float64:
		return uint64(v), true
	case []byte:
		return Uint64(string(v))
	case string:
		v = strings.SplitN(v, ".", 2)[0]
		t, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			return t, true
		} else {
			return 0, false
		}
	default:
		return 0, false
	}
}

func Int(val interface{}) (int, bool) {
	ival, succ := Int64(val)
	return int(ival), succ
}

func Uint(val interface{}) (uint, bool) {
	uval, succ := Uint64(val)
	return uint(uval), succ
}

func Float64(val interface{}) (float64, bool) {
	if val == nil {
		return 0, false
	}
	re_value := reflect.ValueOf(val)
	for re_value.Kind() == reflect.Ptr {
		re_value = re_value.Elem()
		if !re_value.IsValid() {
			return 0, false
		}
		val = re_value.Interface()
		if val == nil {
			return 0, false
		}
		re_value = reflect.ValueOf(val)
	}
	if val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case bool:
		if v {
			return 1, true
		} else {
			return 0, true
		}
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return float64(v), true
	case []byte:
		return Float64(string(v))
	case string:
		t, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return t, true
		} else {
			return 0, false
		}
	default:
		return 0, false
	}
}

func Bool(val interface{}) (bool, bool) {
	ival, succ := Int64(val)
	return ival != 0, succ
}

func IsNil(val interface{}) bool {
	if val == nil {
		return true
	}
	re_value := reflect.ValueOf(val)
	for re_value.Kind() == reflect.Ptr {
		re_value = re_value.Elem()
		if !re_value.IsValid() || re_value.IsNil() {
			return true
		}
		re_value = reflect.ValueOf(re_value.Interface())
	}
	return false
}

func Time(val interface{}) (time.Time, bool) {
	if val == nil {
		return time.Time{}, false
	}
	re_value := reflect.ValueOf(val)
	for re_value.Kind() == reflect.Ptr {
		re_value = re_value.Elem()
		if !re_value.IsValid() {
			return time.Time{}, false
		}
		val = re_value.Interface()
		if val == nil {
			return time.Time{}, false
		}
		re_value = reflect.ValueOf(val)
	}
	if val == nil {
		return time.Time{}, false
	}

	if v, ok := val.(time.Time); ok {
		return v, ok
	} else if v, ok := val.(string); ok {
		tlen := len(v)
		var t time.Time
		var err error
		switch tlen {
		case 8:
			t, err = time.ParseInLocation("20060102", v, time.Local)
		case 10:
			t, err = time.ParseInLocation("2006-01-02", v, time.Local)
		case 19:
			t, err = time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		default:
			return t, false
		}
		if err != nil {
			return t, false
		} else {
			return t, true
		}
	} else {
		return time.Time{}, false
	}
}

func TimePtr(val interface{}) *time.Time {
	t, ok := Time(val)
	if ok {
		return &t
	} else {
		return nil
	}
}
