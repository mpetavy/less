package content

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type buffercontent struct {
	lines []string
}

func NewBufferContent(b []byte) (buffercontent, error) {
	r := bufio.NewReader(bytes.NewReader(b))

	var lines []string
	var line string
	var err error

	for {
		line, err = r.ReadString('\n')

		if err != nil {
			break
		}

		lines = append(lines, line)
	}

	if err == io.EOF {
		return buffercontent{lines}, nil
	} else {
		return buffercontent{}, err
	}
}

func (bc *buffercontent) Line(index int) (string, error) {
	return strings.Trim(bc.lines[index], "\n\r"), nil
}

func (bc *buffercontent) Count() int {
	return len(bc.lines)
}
