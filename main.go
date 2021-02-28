package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mpetavy/less/content"
	"github.com/nsf/termbox-go"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	ErrMissingContent = errors.New("Missing filename (\"less -?\") or piped input")

	help          = flag.Bool("?", false, "Show help")
	filename      = flag.String("f", "", "File you want to parse")
	width, height int
	f             *os.File
	fi            content.Content
	pos           int
	backbuf       []termbox.Cell
	bbw, bbh      int
)

func reset() {
	//termbox.Sync() // cosmestic purpose
}

func bufPrint(x, y int, txt string) {
	for i, r := range txt {
		backbuf[x+i+(y*width)] = termbox.Cell{Ch: r, Fg: termbox.ColorDefault, Bg: termbox.ColorDefault}
	}
}

func readPage(msg string) (int64, error) {
	var format = "%-" + strconv.FormatInt(int64(width), 10) + "s"
	var lines = make([]string, height)
	var show bool
	var err error
	var r int64

	for y := 0; y < len(lines)-1; y++ {
		var line string

		if err == nil {
			line, err = fi.Line(pos + y)
		}

		if err == nil {
			show = true

			r += int64(len(line))

			line = strings.Replace(line, "\t", "    ", -1)

			line = fmt.Sprintf("%3d "+format, y, line)

			if len(line) > width {
				line = line[0:width]
			}
		} else {
			line = fmt.Sprintf(format, " ")
		}

		lines[y] = line
	}

	if show {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.HideCursor()

		if msg == "" {
			msg = ":"
		}

		lines[height-1] = fmt.Sprintf(format, msg)

		for y, v := range lines {
			bufPrint(0, y, v)
		}

		copy(termbox.CellBuffer(), backbuf)

		termbox.SetCursor(len(msg), height-1)

		termbox.Flush()
	}

	return r, err
}

func quit() {
	os.Exit(0)
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func scroll(dir int) {
	v := pos + dir
	v = min(fi.Count()-height, v)
	v = max(0, v)

	if v != pos {
		pos = v

		readPage("")
		reset()
	}
}

func reallocBackBuffer(w, h int) {
	bbw, bbh = w, h
	backbuf = make([]termbox.Cell, w*h)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	width, height = termbox.Size()

	termbox.SetInputMode(termbox.InputEsc)

	reallocBackBuffer(width, height)

	flag.Parse()

	if *filename == "" && len(os.Args) > 1 {
		filename = &os.Args[1]
	}

	if *filename != "" {
		_, err = os.Stat(*filename)
		if err != nil {
			panic(err)
		}

		f, err = os.Open(*filename)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		c, err := content.NewFileContent(f)
		if err != nil {
			panic(err)
		}
		fi = &c
	} else {
		stat, _ := os.Stdin.Stat()

		if (stat.Mode() & os.ModeCharDevice) == 0 {
			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}

			c, err := content.NewBufferContent(b)
			if err != nil {
				panic(err)
			}
			fi = &c
		} else {
			panic(ErrMissingContent)
		}
	}

	pos = 0

	readPage(*filename)

	for {
		switch ev := termbox.PollEvent(); ev.Type {

		case termbox.EventResize:
			reallocBackBuffer(ev.Width, ev.Height)

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				quit()
			case termbox.KeyF1:
				reset()
				fmt.Println("F1 pressed")
			case termbox.KeyF2:
				reset()
				fmt.Println("F2 pressed")
			case termbox.KeyF3:
				reset()
				fmt.Println("F3 pressed")
			case termbox.KeyF4:
				reset()
				fmt.Println("F4 pressed")
			case termbox.KeyF5:
				reset()
				fmt.Println("F5 pressed")
			case termbox.KeyF6:
				reset()
				fmt.Println("F6 pressed")
			case termbox.KeyF7:
				reset()
				fmt.Println("F7 pressed")
			case termbox.KeyF8:
				reset()
				fmt.Println("F8 pressed")
			case termbox.KeyF9:
				reset()
				fmt.Println("F9 pressed")
			case termbox.KeyF10:
				reset()
				fmt.Println("F10 pressed")
			case termbox.KeyF11:
				reset()
				fmt.Println("F11 pressed")
			case termbox.KeyF12:
				reset()
				fmt.Println("F12 pressed")
			case termbox.KeyInsert:
				reset()
				fmt.Println("Insert pressed")
			case termbox.KeyDelete:
				reset()
				fmt.Println("Delete pressed")
			case termbox.KeyHome:
				reset()
				fmt.Println("Home pressed")
			case termbox.KeyEnd:
				reset()
				fmt.Println("End pressed")
			case termbox.KeyPgup:
				scroll(-height)
			case termbox.KeyPgdn:
				scroll(height)
			case termbox.KeyArrowUp:
				scroll(-1)
			case termbox.KeyArrowDown:
				scroll(1)
			case termbox.KeyArrowLeft:
				reset()
				fmt.Println("Arrow Left pressed")
			case termbox.KeyArrowRight:
				reset()
				fmt.Println("Arrow Right pressed")
			case termbox.KeySpace:
				reset()
				readPage("")
			case termbox.KeyBackspace:
				reset()
				fmt.Println("Backspace pressed")
			case termbox.KeyEnter:
				reset()
				fmt.Println("Enter pressed")
			case termbox.KeyTab:
				reset()
				fmt.Println("Tab pressed")
			default:
				// we only want to read a single character or one key pressed event
				reset()
				//fmt.Println("ASCII : ", ev.Ch)
				switch ev.Ch {
				case 'q':
					quit()
				case 'e':
					scroll(-1)
				case 'y':
					scroll(1)
				case 'f':
					scroll(-height)
				case 'b':
					scroll(height)

				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
