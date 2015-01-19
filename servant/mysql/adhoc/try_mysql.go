package main

import (
	"database/sql"
	_ "github.com/funkygao/mysql"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	//DSN = "hellofarm:halfquestfarm4321@tcp(192.168.23.120:3306)/UserShard1?charset=utf8&timeout=10s"
	DSN         = "hellofarm:halfquestfarm4321@tcp(192.168.23.163:3306)/UserShard1?timeout=4s"
	QUERY       = "SELECT * FROM UserInfo WHERE uid=?"
	SCAN_ROWS   = true
	SHOW_ROWS   = false
	USE_PREPARE = true

	DEBUG_ADDR = "127.0.0.1:8765"

	CONN_MAX_IDLE = 5
	CONN_MAX_OPEN = 20

	PARALLAL = 100
	LOOPS    = 1000
)

var (
	wg sync.WaitGroup
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Printf("debug addr: %s/debug/pprof/", DEBUG_ADDR)
	go http.ListenAndServe(DEBUG_ADDR, nil)

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

	t1 := time.Now()
	var rows *sql.Rows
	var stmt *sql.Stmt
	if USE_PREPARE {
		stmt, err = db.Prepare(QUERY)
		if err != nil {
			log.Printf("[%d]: %s", seq, err)
			return
		}

		defer stmt.Close()
	}
	for i := 0; i < LOOPS; i++ {
		if !USE_PREPARE {
			rows, err = db.Query(QUERY, 1)
		} else {
			rows, err = stmt.Query(1)
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
				if ex := rows.Scan(scanArgs...); ex != nil {
					log.Printf("%d[%d]: %s", i+1, seq, ex)
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

			if SHOW_ROWS {
				log.Printf("%d[%d]: %+v", i+1, seq, rowData)
			}

		}

		rows.Close()
	}

	log.Printf("[%3d] %18s %18s", seq, time.Since(t1), time.Since(t1)/LOOPS)
}
