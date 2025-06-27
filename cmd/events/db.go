package events

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (e *EventManager) Connect() {
	new_db, err := sql.Open("sqlite3", e.DbPath)
	if err != nil {
		panic(fmt.Sprintf("critical error in EventManager.Connect(): %v", err))
	}
	fmt.Println("connected to db successfully")
	e.db = new_db
}

func (e *EventManager) CloseConnection() {
	e.db.Close()
	// fmt.Println("closed db connection")
}

func (e *EventManager) CreateTables(query string) {
	_, err := e.db.Exec(query)
	if err != nil {
		panic(fmt.Sprintf("critical error in DBManager.CreateTables(): %v", err))
	}
}

func (e *EventManager) GetEvents(from, until Date, output *map[Date][]Event) {
	query := "SELECT * FROM events WHERE date > ? AND date < ?"

	rows, err := e.db.Query(query, from.Int(), until.Int())
	if err != nil {
		panic(fmt.Sprintf("An error occured, when trying to get rows in EventManager.GetEvents: %v", err))
	}

	toRet := map[Date][]Event{}
	for rows.Next() {
		var event Event
		var date int64
		if err := rows.Scan(&date, &event.Name, &event.Description, &event.Time); err != nil {
			panic(fmt.Sprintf("An error occured when trying to scan rows in EventManager.GetEvents: %v", err))
		}
		t := time.Unix(date, 0)
		toRet[Date{
			Year:  t.Year(),
			Month: t.Month(),
			Day:   t.Day(),
		}] = append(toRet[Date{
			Year:  t.Year(),
			Month: t.Month(),
			Day:   t.Day(),
		}], event)
	}

	*output = toRet
}

var CreateTablesQuery string = `
	CREATE TABLE IF NOT EXISTS events (
		date        INTEGER,
		name        TEXT,
		description TEXT,
		time        INT
	);
`
