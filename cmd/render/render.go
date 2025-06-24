package render

import (
	"fmt"
	"strings"

	"github.com/t0pt/plantica/cmd/events"
	"github.com/t0pt/plantica/cmd/terminal"
)

type Renderer struct {
	Terminal *terminal.TermManager
}

type Column struct {
	Name  string
	Width int
	Cells []Cell
}

type Cell struct {
	Name        string
	Height      int
	Description string
}

func (r *Renderer) RenderSquare() {
	ClearAll()
	fmt.Print(strings.Repeat("-", r.Terminal.Width) + "\n")
	fmt.Print(strings.Repeat("\r|"+strings.Repeat(" ", r.Terminal.Width-2)+"|\n", r.Terminal.Height-3))
	fmt.Print("\r" + strings.Repeat("-", r.Terminal.Width) + "\n")
}

func ClearAll() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func (r *Renderer) RenderColumns(amount int) {
	ClearAll()
	columnWindth := int(r.Terminal.Width / amount)
	column := "|" + strings.Repeat(" ", columnWindth-2) + "|"
	line := "\r" + strings.Repeat(column, amount) + "\n"
	field := strings.Repeat(line, r.Terminal.Height)

	fmt.Print(field)
}

func (r *Renderer) RenderCalendar(days int, focusDate *events.Date) {
	ClearAll()
	field := make([]string, 0, r.Terminal.Height)
	columnWindth := int(r.Terminal.Width / days)
	header := "\r"
	for i := 0; i < days; i++ {
		dateStr := focusDate.AddDays(i - 1).String()
		if len(dateStr) > (columnWindth - 2) {
			dateStr = dateStr[:columnWindth-3]
		}
		header += "|" + SideSpacers(dateStr, columnWindth-2) + "|"
	}
	header += "\n"
	field = append(field, header)
	field = append(field, "\r"+strings.Repeat("—", r.Terminal.Width)+"\n")
	column := "|" + strings.Repeat(" ", columnWindth-2) + "|"
	line := "\r" + strings.Repeat(column, days) + "\n"
	rest := r.Terminal.Height - len(field) - 1
	for i := 0; i < rest; i++ {
		field = append(field, line)
	}
	fmt.Print(strings.Join(field, ""))
}

func SideSpacers(input string, length int) string {
	return strings.Repeat(" ", int((length-len(input))/2)) + input +
		strings.Repeat(" ", int((length-len(input))/2)) +
		strings.Repeat(" ", (length-len(input))%2)
}
