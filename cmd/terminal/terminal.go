package terminal

import (
	"fmt"
	"os"
	"time"

	"github.com/t0pt/plantica/cmd/events"
	"golang.org/x/term"
)

type TermManager struct {
	State      bool // true == raw; false == default
	oldState   *term.State
	Width      int
	Height     int
	FocusDate  *events.Date
	FocusDay   *int
	FocusLine  *int
	MaxRows    *int
	Change     chan bool
	SizeChange chan os.Signal
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
				ter.UpArrow()
			case 66:
				ter.DownArrow()
			case 67:
				ter.RightArrow()
			case 68:
				ter.LeftArrow()
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
	maxRows := 0
	newTerm := TermManager{
		MaxRows: &maxRows,
	}
	newTerm.GetTerminalDimensions()
	go func() {
		for {
			newTerm.GetTerminalDimensions()
			time.Sleep(time.Second)
		}
	}()
	return &newTerm
}

func (ter *TermManager) ChangeSizeDaemon() {
	for {
		<-ter.SizeChange
		newW, newH, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			fmt.Printf("Error getting size: %v\n", err)
			continue
		}
		ter.Width = newW
		ter.Height = newH
		ter.Change <- true
	}
}

func (ter *TermManager) DownArrow() {
	*ter.FocusLine = (*ter.FocusLine + 1)
	if *ter.MaxRows != 0 && *ter.FocusLine > *ter.MaxRows-1 {
		*ter.FocusLine = *ter.MaxRows - 1
	}
	ter.Change <- true
}
func (ter *TermManager) UpArrow() {
	*ter.FocusLine = (*ter.FocusLine - 1)
	if *ter.FocusLine < 0 {
		*ter.FocusLine = 0
	}
	ter.Change <- true
}
func (ter *TermManager) RightArrow() {
	*ter.FocusDay = (*ter.FocusDay + 1)
	if *ter.FocusDay > 4 {
		*ter.FocusDay = 4
		*ter.FocusDate = ter.FocusDate.AddDays(1)
	}
	ter.Change <- true
}
func (ter *TermManager) LeftArrow() {
	*ter.FocusDay = (*ter.FocusDay - 1)
	if *ter.FocusDay < 0 {
		*ter.FocusDay = 0
		*ter.FocusDate = ter.FocusDate.AddDays(-1)
	}
	ter.Change <- true
}
