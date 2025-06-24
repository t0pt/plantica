package main

import (
	"github.com/t0pt/plantica/cmd/events"
	"github.com/t0pt/plantica/cmd/render"
	"github.com/t0pt/plantica/cmd/terminal"
)

func main() {
	mainTerm := terminal.NewTerminal()

	renderer := render.Renderer{
		Terminal: mainTerm,
	}
	today := events.TodayDate()
	renderer.RenderCalendar(5, &today, 0)
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
