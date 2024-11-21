package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ParameterType int

const (
	StringParameter ParameterType = iota
	BooleanParameter
	NumberParameter
	RegExpParameter
	StringListParameter
)

type Parameter struct {
	Description  string        `json:"description"`
	Type         ParameterType `json:"type"`
	DefaultValue string        `json:"default_value"`
}

func ResolveParameters(parameters map[string]Parameter, rawParameters map[string]string) (map[string]any, error) {
	resolvedParameters := make(map[string]any)
	for key, parameter := range parameters {
		var value string
		if v, ok := rawParameters[key]; !ok {
			value = parameter.DefaultValue
		} else {
			value = v
		}
		switch parameter.Type {
		case BooleanParameter:
			if v, err := strconv.ParseBool(value); err != nil {
				return resolvedParameters, fmt.Errorf("can not parse a boolean value from %s: %s", value, err)
			} else {
				resolvedParameters[key] = v
			}
		case NumberParameter:
			if v, err := strconv.ParseInt(value, 10, 64); err != nil {
				return resolvedParameters, fmt.Errorf("can not parse an integer value from %s: %s", value, err)
			} else {
				resolvedParameters[key] = v
			}
		case RegExpParameter:
			if v, err := regexp.Compile(value); err != nil {
				return resolvedParameters, fmt.Errorf("can not parse a regular expression from %s: %s", value, err)
			} else {
				resolvedParameters[key] = v
			}
		case StringListParameter:
			resolvedParameters[key] = strings.Split(value, " ")
		default:
			resolvedParameters[key] = value
		}
	}
	return resolvedParameters, nil
}
