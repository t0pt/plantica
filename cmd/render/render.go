package render

import (
	"fmt"
	"strings"

	"github.com/t0pt/plantica/cmd/events"
	"github.com/t0pt/plantica/cmd/terminal"
)

type Renderer struct {
	Terminal *terminal.TermManager
	Events   *map[events.Date][]events.Event
	EManager *events.EventManager
	lastDate events.Date
}

type Column struct {
	Name  string
	Width int
	Cells []Cell
	Date  events.Date
}

type Cell struct {
	Name        string
	Height      int
	Description string
	Time        int
	IsUsed      bool
}

type LineTime struct {
	Name string
	Time int
}

type TimeLine struct {
	Times []LineTime
	Width int
}

func ClearAll() {
	fmt.Print("\x1b[H\x1b[3J\x1b[2J")
}

// returns celected cell
func (r *Renderer) RenderCalendar(days int, focusDate *events.Date, focusColumn, focusLine int, footer bool) (Cell, events.Date, int) {
	var selectedCell Cell
	var selectedDate events.Date

	dayColumns := []Column{}
	for i := 0; i < days; i++ {
		dayColumns = append(dayColumns,
			Column{
				Name:  focusDate.AddDays(i - 1).String(),
				Cells: []Cell{},
				Date:  focusDate.AddDays(i - 1),
			})
	}

	if r.lastDate != *focusDate {
		r.lastDate = *focusDate
		r.EManager.GetEvents(focusDate.AddDays(-1), focusDate.AddDays(days-1), r.Events)
	}

	for columnIndex, column := range dayColumns {
		eventsMap := *r.Events
		for _, event := range eventsMap[column.Date] {
			dayColumns[columnIndex].Cells = append(dayColumns[columnIndex].Cells, Cell{
				Name:        event.Name,
				Description: event.Description,
				Time:        event.Time,
				IsUsed:      true,
			})
		}
	}

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
		line += "|" + sideSpacers(column.Name, column.Width-2) + "|"
	}
	line += "\n"
	lines = append(lines, line)
	line = "\r" + strings.Repeat("—", r.Terminal.Width) + "\n" // divider
	lines = append(lines, line)

	footerLines := []string{}
	if footer {
		footerLines = r.Footer(2)
	}

	lineBefore := false // if the event before that one has already generated a separation line
	rowCounter := 0     // how many rows with (potential) events are there
	// real business
	for _, lineTime := range timeLine.Times {
		if len(lines) >= r.Terminal.Height-len(footerLines)-1 {
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
				linesInRow = append(linesInRow, "\r|"+sideSpacers(lineTime.Name, timeLine.Width-2)+"|") // starts with linetime
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
						selectedCell = Cell{
							IsUsed: false,
							Time:   lineTime.Time,
						}
						selectedDate = dayColumns[day].Date
						linesInRow[line] += "|" + inverted(strings.Repeat(" ", dayColumns[day].Width-2)) + "|"
					} else {
						linesInRow[line] += "|" + strings.Repeat(" ", dayColumns[day].Width-2) + "|"
					}
					continue
				}
				cell := rowCells[lineTime.Time][day][line]
				if focusLine == rowCounter && focusColumn == day { // +1 because of that later the line before will be printed
					selectedCell = cell
					selectedDate = dayColumns[day].Date
					linesInRow[line] += "|" + inverted(sideSpacers(cell.Name, dayColumns[day].Width-2)) + "|"
				} else {
					linesInRow[line] += "|" + sideSpacers(cell.Name, dayColumns[day].Width-2) + "|"
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
	for len(lines) < r.Terminal.Height-len(footerLines)-1 {
		line := "\r" + strings.Repeat(" ", timeLine.Width)
		for day, column := range dayColumns {
			if focusLine == rowCounter && focusColumn == day {
				selectedCell = Cell{
					IsUsed: false,
					Time:   -1,
				}
				selectedDate = dayColumns[day].Date
				line += "|" + inverted(strings.Repeat(" ", column.Width-2)) + "|"
			} else {
				line += "|" + strings.Repeat(" ", column.Width-2) + "|"
			}
		}
		line += "\n"
		lines = append(lines, line)
		rowCounter++
	}
	ClearAll()
	fmt.Print(strings.Join(lines, ""))
	fmt.Print(strings.Join(footerLines, ""))
	return selectedCell, selectedDate, rowCounter
}

// space represents the minimum amount of space between footer items
func (r *Renderer) Footer(space int) []string {
	lines := make([]string, 0, 2)
	line := []string{}

	length := 0
	for _, content := range footerContent {
		if length+len(content) < r.Terminal.Width { // the line still has free space
			line = append(line, content)
			length += len(content) + space
		} else { // new line should be created
			length -= len(line) * space
			freeSpace := r.Terminal.Width - length
			if len(line) > 1 {
				lines = append(lines, "\r"+strings.Join(line, strings.Repeat(" ", int(freeSpace/(len(line)-1))))+"\n")
			} else {
				lines = append(lines, "\r"+line[0]+"\n")
			}
			line = []string{content}
			length = len(content) + space
		}
	}
	if len(line) != 0 {
		length -= len(line) * space
		freeSpace := r.Terminal.Width - length
		if len(line) > 1 {
			lines = append(lines, "\r"+strings.Join(line, strings.Repeat(" ", int(freeSpace/(len(line)-1))))+"\n")
		} else {
			lines = append(lines, "\r"+line[0]+"\n")
		}
	}
	return lines
}

func sideSpacers(input string, length int) string {
	if len(input) >= length {
		return input[:length]
	}
	return strings.Repeat(" ", int((length-len(input))/2)) + input +
		strings.Repeat(" ", int((length-len(input))/2)) +
		strings.Repeat(" ", (length-len(input))%2)
}

func inverted(input string) string {
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

var footerContent []string = []string{
	"Navigate-" + inverted("Arrows"),
	"New event-" + inverted("ENTER"),
	"Exit-" + inverted("Ctrl+C"),
}
