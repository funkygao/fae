package main

import (
	"database/sql"
	_ "github.com/funkygao/mysql"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	//DSN = "hellofarm:halfquestfarm4321@tcp(192.168.23.120:3306)/UserShard1?charset=utf8&timeout=10s"
	DSN         = "hellofarm:halfquestfarm4321@tcp(192.168.23.163:3306)/UserShard1?timeout=4s"
	QUERY       = "SELECT * FROM UserInfo WHERE uid=?"
	SCAN_ROWS   = false
	USE_PREPARE = false

	CONN_MAX_IDLE = 5
	CONN_MAX_OPEN = 20

	PARALLAL = 100
)

var (
	wg sync.WaitGroup
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	t1 := time.Now()
	for i := 0; i < PARALLAL; i++ {
		wg.Add(1)
		go runDb(i)
	}

	wg.Wait()

	log.Println()
	log.Println(time.Since(t1))
}

func runDb(seq int) {
	defer wg.Done()

	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	if CONN_MAX_IDLE > 0 {
		db.SetMaxIdleConns(CONN_MAX_IDLE)
	}
	if CONN_MAX_OPEN > 0 {
		db.SetMaxOpenConns(CONN_MAX_OPEN)
	}

	const N = 5000
	t1 := time.Now()
	var rows *sql.Rows
	for i := 0; i < N; i++ {
		if !USE_PREPARE {
			rows, err = db.Query(QUERY, 1)
		} else {
			stmt, e := db.Prepare(QUERY)
			if e != nil {
				log.Printf("%d[%d]: %s", i+1, seq, e)
				return
			}
			rows, err = stmt.Query(1)
			defer stmt.Close()
		}

		if err != nil {
			log.Printf("%d[%d]: %s", i+1, seq, err)
			return
		}

		if SCAN_ROWS {
			cols, err := rows.Columns()
			if err != nil {
				log.Printf("%d[%d]: %s", i+1, seq, err)
				return
			}

			rowData := make([][]string, 0)
			for rows.Next() {
				rowValues := make([]string, len(cols))
				rawRowValues := make([]sql.RawBytes, len(cols))
				scanArgs := make([]interface{}, len(cols))
				for i, _ := range cols {
					scanArgs[i] = &rawRowValues[i]
				}
				if appErr := rows.Scan(scanArgs...); appErr != nil {
					log.Printf("%d[%d]: %s", i+1, seq, appErr)
					return
				}

				for i, raw := range rawRowValues {
					if raw == nil {
						rowValues[i] = "NULL"
					} else {
						rowValues[i] = string(raw)
					}
				}

				rowData = append(rowData, rowValues)
			}

			if rows.Err() != nil {
				log.Printf("%d[%d]: %s", i+1, seq, rows.Err())
				return
			}

			//log.Printf("%+v", rowData)
		}

		rows.Close()
	}

	log.Println(seq, time.Since(t1), time.Since(t1)/N)
}
