package endpoints

import (
	"fmt"
	"github.com/emersion/go-ical"
	"github.com/gin-gonic/gin"
	"icarusserver/internal/processors"
	"io"
	"strings"
)

var supportedProcessors = []processors.ProcessorSupport{
	processors.FilterSupport{},
}

type ProcessorsEndpoint struct{}

func (p ProcessorsEndpoint) Register(g *gin.RouterGroup) error {
	g.GET("/processors", p.getProcessors)
	g.POST("/processors/:processor", p.process)
	return nil
}

func (p ProcessorsEndpoint) getProcessors(c *gin.Context) {
	processorMap := p.getProcessorMap()
	c.JSON(200, processorMap)
}

func (p ProcessorsEndpoint) getProcessorMap() map[string]processors.ProcessorSupportMetadata {
	processorMap := make(map[string]processors.ProcessorSupportMetadata)
	for _, processor := range supportedProcessors {
		m := processor.GetMetadata()
		processorMap[m.EndpointKey] = m
	}
	return processorMap
}

func (p ProcessorsEndpoint) process(context *gin.Context) {
	var b struct {
		Calendar string `json:"calendar"`
	}
	var input io.Reader
	if err := context.BindJSON(&b); err != nil {
		if f, err := context.FormFile("file"); err != nil {
			context.JSON(400, gin.H{
				"error": "Calendar data needs to be specified using a JSON object with the key 'calendar' or as an upload with a 'file' form field",
			})
			return
		} else {
			if fh, err := f.Open(); err != nil {
				context.JSON(500, gin.H{
					"error": "Can't read uploaded calendar",
				})
				return
			} else {
				input = fh
			}
		}
	} else {
		input = strings.NewReader(b.Calendar)
	}
	var inputCalendar ical.Calendar
	if i, err := ical.NewDecoder(input).Decode(); err != nil {
		context.JSON(400, gin.H{
			"error": "Input is not valid ICS data",
		})
		return
	} else {
		inputCalendar = *i
	}

	var outputCalendar = ical.NewCalendar()
	outputCalendar.Props = inputCalendar.Props

	for _, processor := range supportedProcessors {
		m := processor.GetMetadata()
		if m.EndpointKey == context.Param("processor") {
			parameters := make(map[string]string)
			for key := range m.Parameters {
				if v, ok := context.GetQuery(key); ok {
					parameters[key] = v
				}
			}
			if err := processor.Process(parameters, inputCalendar, outputCalendar); err != nil {
				context.JSON(500, gin.H{
					"error": fmt.Sprintf("Error running processor: %s", err),
				})
				return
			}
			o := strings.Builder{}
			if err := ical.NewEncoder(&o).Encode(outputCalendar); err != nil {
				context.JSON(500, gin.H{
					"error": fmt.Sprintf("Error converting calendar to ics: %s", err),
				})
				return
			}
			context.JSON(200, gin.H{
				"calendar": o.String(),
			})
		}
	}
}

var _ Endpoint = ProcessorsEndpoint{}
