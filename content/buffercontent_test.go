package content

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBufferContent(t *testing.T) {
	lines := []string{"1", "2", "3"}

	var b bytes.Buffer

	for _, l := range lines {
		b.Write([]byte(l + "\n"))
	}

	c, err := NewBufferContent(b.Bytes())
	if err != nil {
		panic(err)
	}

	checker := func(lines []string, fi Content) {
		for i := range lines {
			l, err := fi.Line(i)
			if err != nil {
				panic(err)
			}

			if l != lines[i] {
				panic(fmt.Errorf("Not equal. Value: %v, Expected: %v", l, lines[i]))
			}
		}
	}

	checker(lines, &c)
}
