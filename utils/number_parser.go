package utils

import (
	"encoding/json"
	"strconv"
)

// ParseInt parse interface to int
func ParseInt(value interface{}, args ...int) int {
	if len(args) > 0 {
		return int(ParseInt64(value, int64(args[0])))
	}

	return int(ParseInt64(value))
}

// ParseInt64 parse interface to int64
func ParseInt64(value interface{}, args ...int64) int64 {
	switch n := value.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return n
	case float32:
		return int64(n)
	case float64:
		return int64(n)
	case uint:
		return int64(n)
	case uint8:
		return int64(n)
	case uint16:
		return int64(n)
	case uint32:
		return int64(n)
	case uint64:
		return int64(n)
	case bool:
		if n {
			return 1
		}
		return 0
	case json.Number:
		if valueInt, err := value.(json.Number).Int64(); err == nil {
			return valueInt
		}
	case string:
		if valueInt, err := strconv.ParseInt(value.(string), 10, 64); err == nil {
			return valueInt
		}
	}

	if len(args) > 0 {
		return args[0]
	}

	return -1
}

// ParseFloat64 parse interface to float64
func ParseFloat64(value interface{}, args ...float64) float64 {
	switch n := value.(type) {
	case int:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	case uint:
		return float64(n)
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case bool:
		if n {
			return 1
		}
		return 0
	case json.Number:
		if valueFloat64, err := value.(json.Number).Float64(); err == nil {
			return valueFloat64
		}
	case string:
		if valueFloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
			return valueFloat64
		}
	}

	if len(args) > 0 {
		return args[0]
	}

	return -1
}
