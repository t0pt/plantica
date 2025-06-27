package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/t0pt/plantica/cmd/events"
	"github.com/t0pt/plantica/cmd/render"
	"github.com/t0pt/plantica/cmd/terminal"
)

func main() {
	mainTerm := terminal.NewTerminal()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGWINCH)
	mainTerm.SizeChange = sigc
	go mainTerm.ChangeSizeDaemon()

	renderer := render.Renderer{
		Terminal: mainTerm,
	}
	focusDay := 1
	focusLine := 5
	focusDate := events.TodayDate().AddDays(0)
	mainTerm.FocusDay = &focusDay
	mainTerm.FocusLine = &focusLine
	change := make(chan bool)
	mainTerm.Change = change
	go func() {
		for {
			renderer.RenderCalendar(5, &focusDate, focusDay, focusLine)
			<-change
		}
	}()
	mainTerm.EnableRaw()
	mainTerm.Listen()
	mainTerm.DisableRaw()

	// fmt.Print("Press arrow keys. Press Ctrl+C to exit.\r\n")
	// fmt.Print("Enter your name: ")
	// reader := bufio.NewReader(os.Stdin)
	// name, _ := reader.ReadString('\n')
	// fmt.Println(name)
	// mainTerm.EnableRaw()
	// mainTerm.Speak()

	// fmt.Println("exiting")

	// mainTerm.DisableRaw()
}
