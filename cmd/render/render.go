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

func (r *Renderer) RenderCalendar(days int, focusDate *events.Date, focusRow int) {
	ClearAll()
	field := []Column{}

	timeLine := Column{
		Width: 9,
		Cells: generateTimeLine(),
	}

	dayColumns := []Column{}
	for i := 0; i < days; i++ {
		dayColumns = append(dayColumns,
			Column{
				Name:  focusDate.AddDays(i - 1).String(),
				Cells: []Cell{},
			})
	}

	field = append(field, timeLine)
	field = append(field, dayColumns...)

	dedicatedWidth := 0 // distribue free space between columns without width
	notDedicatedIds := []int{}
	for i := 0; i < len(field); i++ {
		if field[i].Width > 2 {
			dedicatedWidth += field[i].Width
		} else {
			notDedicatedIds = append(notDedicatedIds, i)
		}
	}
	spacePerColumn := (r.Terminal.Width - dedicatedWidth) / len(notDedicatedIds)
	for i := 0; i < len(notDedicatedIds); i++ {
		field[notDedicatedIds[i]].Width = spacePerColumn
	}

	cellMatrix := [][]Cell{}
	heightRow := map[int]int{} // contains height of each row
	for i := 0; i < len(field); i++ {
		cellMatrix = append(cellMatrix, field[i].Cells)
		for j := 0; j < len(field[i].Cells); j++ {
			if heightRow[j] < field[i].Cells[j].Height {
				heightRow[j] = field[i].Cells[j].Height
			}
		}
	}

	sumHeights := 0
	currentRow := 0
	lines := make([]string, r.Terminal.Height-1)
	for i := 0; i < len(lines); i++ {
		lines[i] = "\r"
	}
	for lin := 0; lin < len(lines); lin++ {
		newRow := false
		divider := false
		if lin-1 == sumHeights { // divider
			divider = true
		} else if lin-1 > sumHeights { // start new row
			newRow = true
			currentRow += 1
			sumHeights += heightRow[currentRow]
		}
		for col := 0; col < len(field); col++ {
			if lin == 0 { // header
				lines[lin] += "|" + SideSpacers(field[col].Name, field[col].Width-2) + "|"
			} else if lin == 1 { // divider
				lines[lin] += strings.Repeat("—", field[col].Width)
			} else { // body
				if divider { // divider
					lines[lin] += strings.Repeat("—", field[col].Width)
				} else if newRow { // start new row
					if len(cellMatrix) >= col && len(cellMatrix[col]) >= currentRow {
						lines[lin] += "|" + SideSpacers(cellMatrix[col][currentRow-1].Name, field[col].Width-2) + "|"
					} else {
						lines[lin] += "|" + strings.Repeat(" ", field[col].Width-2) + "|"
					}
				} else {
					lines[lin] += "|" + strings.Repeat(" ", field[col].Width-2) + "|"
				}
			}
		}
	}
	for i := 0; i < len(lines); i++ {
		lines[i] += "\n"
	}

	fmt.Print(strings.Join(lines, ""))
}

func SideSpacers(input string, length int) string {
	if len(input) >= length {
		return input[:length]
	}
	return strings.Repeat(" ", int((length-len(input))/2)) + input +
		strings.Repeat(" ", int((length-len(input))/2)) +
		strings.Repeat(" ", (length-len(input))%2)
}

func generateTimeLine() []Cell {
	toRet := []Cell{}
	for i := 6; i != 5; i++ {
		if i == 24 {
			i = 0
		}
		hours := ""
		if i < 10 {
			hours = fmt.Sprintf("0%d", i)
		} else {
			hours = fmt.Sprintf("%d", i)
		}
		toRet = append(toRet,
			Cell{
				Name:   hours + ":00",
				Height: 3,
			})
	}
	return toRet
}
