package utils

import (
	"github.com/json-iterator/go"
	"github.com/yiplee/structs"
)

func newStrcuts(object interface{}) *structs.Struct {
	s := structs.New(object)
	s.FlattenAnonymous = true
	s.TagName = "json"
	return s
}

// SelectFields output object to map
func SelectFields(object interface{}, fields ...string) map[string]interface{} {
	s := newStrcuts(object)
	sMap := s.Map()

	out := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		if object, ok := sMap[field]; ok {
			out[field] = object
		}
	}

	return out
}

// UnselectFields output object to map
func UnselectFields(object interface{}, fields ...string) map[string]interface{} {
	s := newStrcuts(object)
	sMap := s.Map()

	for _, field := range fields {
		delete(sMap, field)
	}

	return sMap
}

// Map generate map
func Map(objects ...interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for i := 0; i < len(objects)-1; i += 2 {
		key := objects[i].(string)
		value := objects[i+1]
		result[key] = value
	}

	return result
}

// JSONString generate JSON
func JSONString(objects ...interface{}) (string, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(Map(objects...))
	return string(data), err
}
