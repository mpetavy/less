package content

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileContent(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
		err = os.Remove(f.Name())
		if err != nil {
			panic(err)
		}
	}()

	lines := []string{"1", "2", "3"}

	f.WriteString(lines[0] + "\n")
	f.WriteString(lines[1] + "\n")
	f.WriteString(lines[2] + "\n")

	c, err := NewFileContent(f)
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

	lines = append(lines, "a", "b", "c")

	f.WriteString(lines[3] + "\n")
	f.WriteString(lines[4] + "\n")
	f.WriteString(lines[5] + "\n")

	c.Update()

	checker(lines, &c)
}
