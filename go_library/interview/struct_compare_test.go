package interview

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Stu struct {
	Name string
	Age  int
}

func TestStruct(t *testing.T) {
	stru1 := Stu{"w", 18}
	stru2 := Stu{"w", 18}

	assert.Equal(t, stru1, stru2, "should be same")
}
