package processors

import (
	"fmt"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
	"icarusserver/internal"
	"regexp"
)

type FilterSupport struct{}

func (FilterSupport) GetMetadata() ProcessorSupportMetadata {
	return ProcessorSupportMetadata{
		EndpointKey: "filter",
		Name:        "Filter",
		Description: "Filters events based on field values",
		Parameters: map[string]internal.Parameter{
			"inverse": {
				Type:         internal.BooleanParameter,
				DefaultValue: "false",
				Description:  "Whether to invert the filter",
			},
			"selector": {
				Type:         internal.RegExpParameter,
				DefaultValue: "",
				Description:  "The value to search for",
			},
			"selectorFields": {
				Type:         internal.StringListParameter,
				DefaultValue: "SUMMARY DESCRIPTION",
				Description:  "Whitespace separated list of event fields to filter on",
			},
		},
	}
}

func (f FilterSupport) Process(parameters map[string]string, input ical.Calendar, output *ical.Calendar) error {
	var inverse bool
	var selector *regexp.Regexp
	var selectorFields []string

	if par, err := internal.ResolveParameters(f.GetMetadata().Parameters, parameters); err != nil {
		return fmt.Errorf("error parsing filter parameters: %s", err)
	} else {
		inverse = par["inverse"].(bool)
		selector = par["selector"].(*regexp.Regexp)
		selectorFields = par["selectorFields"].([]string)
	}

	p := processors.FilterProcessor{
		Inverse: inverse,
	}
	t := processors.Toolbox{
		TextSelectorPattern: selector,
		TextSelectorProps:   selectorFields,
	}
	p.SetToolbox(t)
	return p.Process(input, output)
}
