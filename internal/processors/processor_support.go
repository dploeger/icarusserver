package processors

import (
	"github.com/emersion/go-ical"
	"icarusserver/internal"
)

type ProcessorSupportMetadata struct {
	EndpointKey string                        `json:"endpoint_key"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	Parameters  map[string]internal.Parameter `json:"parameters"`
}

type ProcessorSupport interface {
	GetMetadata() ProcessorSupportMetadata
	Process(parameters map[string]string, input ical.Calendar, output *ical.Calendar) error
}
