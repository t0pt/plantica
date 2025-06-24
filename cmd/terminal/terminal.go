package terminal

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

type TermManager struct {
	State    bool // true == raw; false == default
	oldState *term.State
	Width    int
	Height   int
}

func (ter *TermManager) Speak() {
	buf := make([]byte, 3)
	for {
		os.Stdin.Read(buf)
		if buf[0] == 27 && buf[1] == 91 {
			switch buf[2] {
			case 65:
				fmt.Print("↑ Up arrow\r\n")
			case 66:
				fmt.Print("↓ Down arrow\r\n")
			case 67:
				fmt.Print("→ Right arrow\r\n")
			case 68:
				fmt.Print("← Left arrow\r\n")
			}
		} else {
			if buf[0] == 27 && buf[1] == 0 && buf[2] == 0 { // ESC
				fmt.Printf("You pressed: ESC\r\n")
			} else if buf[1] == 0 && buf[2] == 0 { // letters and everything else, ignore F keys
				switch buf[0] {
				case 127:
					fmt.Print("You pressed: BACKSPACE\r\n")
				case 3: // thats ctr+c
					return
				default:
					fmt.Printf("You pressed: %q\r\n", buf[0])
				}
			} else {
				fmt.Println(buf)
			}
		}
		buf = make([]byte, 3)
	}
}

func (ter *TermManager) Listen() {
	buf := make([]byte, 3)
	for {
		os.Stdin.Read(buf)
		if buf[0] == 27 && buf[1] == 91 {
			switch buf[2] {
			case 65:
			case 66:
			case 67:
			case 68:
			}
		} else {
			if buf[0] == 27 && buf[1] == 0 && buf[2] == 0 { // ESC
			} else if buf[1] == 0 && buf[2] == 0 { // letters and everything else, ignore F keys
				switch buf[0] {
				case 127:
				case 3: // thats ctr+c
					return
				default:
				}
			} else {
			}
		}
		buf = make([]byte, 3)
	}
}

func (ter *TermManager) EnableRaw() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error setting raw mode:", err)
		os.Exit(1)
	}
	ter.oldState = oldState
	ter.State = true
}

func (ter *TermManager) DisableRaw() {
	term.Restore(int(os.Stdin.Fd()), ter.oldState)
	ter.State = false
}

func (ter *TermManager) GetTerminalDimensions() {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting terminal dimensions:", err)
		os.Exit(1)
	}
	ter.Width = width
	ter.Height = height
}

func NewTerminal() *TermManager {
	newTerm := TermManager{}
	newTerm.GetTerminalDimensions()
	go func() {
		for {
			newTerm.GetTerminalDimensions()
			time.Sleep(time.Second)
		}
	}()
	return &newTerm
}
