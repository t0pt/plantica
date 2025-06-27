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
	Day   int
}

type Cell struct {
	Name        string
	Height      int
	Description string
	Time        int
}

type LineTime struct {
	Name string
	Time int
}

type TimeLine struct {
	Times []LineTime
	Width int
}

func (r *Renderer) RenderSquare() {
	ClearAll()
	fmt.Print(strings.Repeat("-", r.Terminal.Width) + "\n")
	fmt.Print(strings.Repeat("\r|"+strings.Repeat(" ", r.Terminal.Width-2)+"|\n", r.Terminal.Height-3))
	fmt.Print("\r" + strings.Repeat("-", r.Terminal.Width) + "\n")
}

func ClearAll() {
	fmt.Print("\x1b[H\x1b[3J\x1b[2J")
}

func (r *Renderer) RenderColumns(amount int) {
	ClearAll()
	columnWindth := int(r.Terminal.Width / amount)
	column := "|" + strings.Repeat(" ", columnWindth-2) + "|"
	line := "\r" + strings.Repeat(column, amount) + "\n"
	field := strings.Repeat(line, r.Terminal.Height)

	fmt.Print(field)
}

// returns celected cell
func (r *Renderer) RenderCalendar(days int, focusDate *events.Date, focusColumn, focusLine int) *Cell {
	ClearAll()
	var selectedCell *Cell

	dayColumns := []Column{}
	for i := 0; i < days; i++ {
		dayColumns = append(dayColumns,
			Column{
				Name:  focusDate.AddDays(i - 1).String(),
				Cells: []Cell{},
				Day:   focusDate.AddDays(i - 1).Day,
			})
	}
	dayColumns[1].Cells = append(dayColumns[1].Cells, Cell{Name: "skibidi", Time: 8})
	dayColumns[1].Cells = append(dayColumns[1].Cells, Cell{Name: "skibidii", Time: 8})
	dayColumns[1].Cells = append(dayColumns[1].Cells, Cell{Name: "pididi", Time: 8})
	dayColumns[1].Cells = append(dayColumns[1].Cells, Cell{Name: "skibidi", Time: 3})
	dayColumns[3].Cells = append(dayColumns[3].Cells, Cell{Name: "skibidi", Time: 10})
	dayColumns[4].Cells = append(dayColumns[4].Cells, Cell{Name: "skibidi", Time: 12})

	timeLine := generateTimeLine(6, 20)
	dedicatedWidth := timeLine.Width // distribue free space between columns without width
	notDedicatedIds := []int{}
	for i := 0; i < len(dayColumns); i++ {
		if dayColumns[i].Width != 0 {
			dedicatedWidth += dayColumns[i].Width
		} else {
			notDedicatedIds = append(notDedicatedIds, i)
		}
	}
	spacePerColumn := (r.Terminal.Width - dedicatedWidth) / len(notDedicatedIds)
	for i := 0; i < len(notDedicatedIds); i++ {
		dayColumns[notDedicatedIds[i]].Width = spacePerColumn
	}

	rowCells := map[int][][]Cell{} // map[time][column][index in this time] = Cell
	for i := 6; len(rowCells) < 20; i++ {
		if i == 24 {
			i = 0
		}
		rowCells[i] = [][]Cell{}
		for j := 0; j < 5; j++ {
			rowCells[i] = append(rowCells[i], []Cell{})
		}
	}
	for j := 0; j < 5; j++ {
		rowCells[-1] = append(rowCells[-1], []Cell{})
	}
	for day, column := range dayColumns {
		for _, cell := range column.Cells {
			if cell.Time < 6 && cell.Time > 2 || cell.Time > 23 { // if the time in cell is not to be displayed => to the pending
				cell.Time = -1
			}
			rowCells[cell.Time][day] = append(rowCells[cell.Time][day], cell)
		}
	}

	lines := []string{}
	line := "\r|" + strings.Repeat(" ", timeLine.Width-2) + "|" // header
	for _, column := range dayColumns {
		line += "|" + SideSpacers(column.Name, column.Width-2) + "|"
	}
	line += "\n"
	lines = append(lines, line)
	line = "\r" + strings.Repeat("—", r.Terminal.Width) + "\n" // divider
	lines = append(lines, line)

	lineBefore := false
	rowCounter := 0
	// real business
	for _, lineTime := range timeLine.Times {
		if len(lines) >= r.Terminal.Height-1 {
			break
		}
		maxLines := 1 // max events in one time
		eventsExist := false
		for _, cells := range rowCells[lineTime.Time] {
			if len(cells) > maxLines {
				maxLines = len(cells)
			}
			if len(cells) > 0 {
				eventsExist = true
			}
		}
		linesInRow := make([]string, 0, maxLines)
		for i := range maxLines {
			if i == 0 {
				linesInRow = append(linesInRow, "\r|"+SideSpacers(lineTime.Name, timeLine.Width-2)+"|") // starts with linetime
			} else {
				linesInRow = append(linesInRow, "\r"+strings.Repeat(" ", timeLine.Width))
			}
		}
		for line := range maxLines {
			if _, ok := rowCells[lineTime.Time]; !ok { // no events at this time at all
				for day := range len(dayColumns) {
					linesInRow[line] += "|" + strings.Repeat(" ", dayColumns[day].Width-2) + "|"
				}
				continue
			}
			for day := range len(dayColumns) {
				if len(rowCells[lineTime.Time][day])-1 < line { // nothing in that time in that day in that line
					if focusLine == rowCounter && focusColumn == day { // +1 because of that later the line before will be printed
						linesInRow[line] += "|" + Inverted(strings.Repeat(" ", dayColumns[day].Width-2)) + "|"
					} else {
						linesInRow[line] += "|" + strings.Repeat(" ", dayColumns[day].Width-2) + "|"
					}
					continue
				}
				cell := rowCells[lineTime.Time][day][line]
				if focusLine == rowCounter && focusColumn == day { // +1 because of that later the line before will be printed
					selectedCell = &cell
					linesInRow[line] += "|" + Inverted(SideSpacers(cell.Name, dayColumns[day].Width-2)) + "|"
				} else {
					linesInRow[line] += "|" + SideSpacers(cell.Name, dayColumns[day].Width-2) + "|"
				}
			}
			rowCounter += 1
		}
		for i := range maxLines {
			linesInRow[i] += "\n" // ends each line
		}
		if !lineBefore && eventsExist {
			linesInRow = append([]string{"\r" + strings.Repeat("—", r.Terminal.Width) + "\n"}, linesInRow...)
		}
		if eventsExist && lineTime.Time != -1 {
			linesInRow = append(linesInRow, "\r"+strings.Repeat("—", r.Terminal.Width)+"\n")
			lineBefore = true
		} else {
			lineBefore = false
		}
		lines = append(lines, linesInRow...)
	}
	fmt.Print(strings.Join(lines, ""))
	return selectedCell
}

func SideSpacers(input string, length int) string {
	if len(input) >= length {
		return input[:length]
	}
	return strings.Repeat(" ", int((length-len(input))/2)) + input +
		strings.Repeat(" ", int((length-len(input))/2)) +
		strings.Repeat(" ", (length-len(input))%2)
}

func Inverted(input string) string {
	return "\x1b[7m" + input + "\x1b[0m"
}

func generateTimeLine(from, amount int) TimeLine {
	toRet := TimeLine{
		Width: 9,
		Times: []LineTime{},
	}
	for i := from; len(toRet.Times) < amount; i++ {
		if i == 24 {
			i = 0
		}
		hours := ""
		if i < 10 {
			hours = fmt.Sprintf("0%d", i)
		} else {
			hours = fmt.Sprintf("%d", i)
		}
		toRet.Times = append(toRet.Times, LineTime{
			Name: hours + ":00",
			Time: i,
		})
	}
	toRet.Times = append(toRet.Times, LineTime{
		Name: "pending",
		Time: -1,
	})
	return toRet
}
