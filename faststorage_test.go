package faststorage

import (
	"github.com/gomodule/redigo/internal/redistest"
	"testing"
)

func TestNew(t *testing.T){
	redistest.Dial()
}
