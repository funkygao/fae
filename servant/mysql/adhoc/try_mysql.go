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

var wg sync.WaitGroup

func init() {
	runtime.GOMAXPROCS(8)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	t1 := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go runDb(i)
	}

	wg.Wait()

	log.Println()
	log.Println(time.Since(t1))
}

func runDb(no int) {
	defer wg.Done()

	db, err := sql.Open("mysql", "hellofarm:halfquestfarm4321@tcp(192.168.23.163:3306)/UserShard1?timeout=4s")
	//db, err := sql.Open("mysql", "hellofarm:halfquestfarm4321@tcp(192.168.23.120:3306)/UserShard1?charset=utf8&timeout=10s")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	//db.SetMaxIdleConns(5)
	//db.SetMaxOpenConns(20)

	query := "SELECT * FROM UserInfo WHERE uid=?"
	const N = 5000
	t1 := time.Now()
	for i := 0; i < N; i++ {
		rows, err := db.Query(query, 1)
		if err != nil {
			log.Printf("%d[%d]: %s", i+1, no, err)
			return
		}

		if false {
			cols, err := rows.Columns()
			if err != nil {
				log.Printf("%d: %s", i+1, err)
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
					log.Printf("%s", appErr)
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
				log.Println(rows.Err())
				return
			}

			//log.Printf("%+v", rowData)
		}

		rows.Close()
	}

	db.Close()
	log.Println(no, time.Since(t1), time.Since(t1)/N)
}
