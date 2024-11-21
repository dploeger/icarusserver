package processors

import (
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFilterSupport_Process(t *testing.T) {
	f, _ := os.Open("../../tests/example.ics")
	c, _ := ical.NewDecoder(f).Decode()
	s := FilterSupport{}
	o := ical.NewCalendar()
	if assert.NoError(t, s.Process(map[string]string{}, *c, o)) {
		assert.Len(t, o.Children, 2)
	}
}
