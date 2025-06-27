package main

import (
	"fmt"
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

	focusDay := 1
	focusLine := 5
	focusDate := events.TodayDate().AddDays(0)
	mainTerm.FocusDay = &focusDay
	mainTerm.FocusLine = &focusLine

	change := make(chan bool)
	mainTerm.Change = change

	var Events = map[events.Date][]events.Event{}
	eventManager := events.EventManager{
		DbPath: "./plantica.db",
	}
	eventManager.Connect()
	defer eventManager.CloseConnection()
	eventManager.GetEvents(events.TodayDate().AddDays(-1), events.TodayDate().AddDays(3), &Events)

	renderer := render.Renderer{
		Terminal: mainTerm,
		Events:   &Events,
	}
	go func() {
		for {
			cell, date, maxRows := renderer.RenderCalendar(5, &focusDate, focusDay, focusLine)
			*mainTerm.MaxRows = maxRows
			fmt.Print("\r", cell, "\n")
			fmt.Print("\r", date, "\n")
			<-change
		}
	}()
	mainTerm.EnableRaw()
	mainTerm.Listen()
	render.ClearAll()
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
