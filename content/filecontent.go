package content

import (
	"io"
	"os"
	"strings"
)

type filecontent struct {
	fileLength int64
	file       *os.File
	lines      []int64
}

func NewFileContent(f *os.File) (filecontent, error) {
	fi := filecontent{file: f}

	err := fi.Update()
	if err != nil {
		return filecontent{}, err
	}

	return fi, nil
}

func (fi *filecontent) Update() error {
	var pos int64

	if len(fi.lines) == 0 {
		fstat, err := os.Stat(fi.file.Name())
		if err != nil {
			return err
		}

		fi.fileLength = fstat.Size()
		fi.lines = make([]int64, 0, 100)

		pos = 0

		_, err = fi.file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
	} else {
		pos = fi.fileLength

		fstat, err := os.Stat(fi.file.Name())
		if err != nil {
			return err
		}

		fi.fileLength = fstat.Size()

		if len(fi.lines) > 0 {
			fi.lines = fi.lines[0 : len(fi.lines)-1]
		}

		_, err = fi.file.Seek(pos, io.SeekStart)
		if err != nil {
			return err
		}
	}

	newLine := true

	for {
		if newLine {
			fi.lines = append(fi.lines, pos)
		}

		b := make([]byte, 1)

		_, err := fi.file.Read(b)
		if err == io.EOF {
			break
		} else {
			if err != nil {
				return err
			}
		}

		newLine = b[0] == '\n'
		pos++
	}

	return nil
}

func (fi *filecontent) seekline(index int) (int64, error) {
	if index < len(fi.lines) {
		return fi.lines[index], nil
	} else {
		return -1, io.EOF
	}
}

func (fi *filecontent) Line(index int) (string, error) {
	sp, err := fi.seekline(index)
	if err != nil {
		return "", err
	}

	_, err = fi.file.Seek(sp, io.SeekStart)
	if err != nil {
		return "", err
	}

	var l int64
	if index < (len(fi.lines) - 1) {
		l = fi.lines[index+1] - fi.lines[index]
	} else {
		l = fi.fileLength - fi.lines[index]
	}

	b := make([]byte, l)

	_, err = fi.file.Read(b)
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.Trim(string(b), "\n\r"), nil
}

func (fi *filecontent) Count() int {
	return len(fi.lines)
}
