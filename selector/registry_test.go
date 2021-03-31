package selector_test

import (
	. "github.com/bukowa/micro/selector"
	"testing"
)

func TestBaseRegistry_Register(t *testing.T) {
	reg := &BaseRegistry{}
	reg.Register("1", NewBaseSelector("2", func(i interface{}) []Score {
		return nil
	}))
	if len(reg.Scored("")) > 0 {
		t.Error()
	}
}

