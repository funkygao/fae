package main

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"labix.org/v2/mgo/bson"
	"log"
	"math/rand"
	"sync"
	"time"
)

func runSession(proxy *proxy.Proxy, wg *sync.WaitGroup, round int, seq int) {
	report.updateConcurrency(1)
	report.incSessions()
	defer func() {
		wg.Done()
		report.updateConcurrency(-1)
		//log.Printf("session{round^%d seq^%d} done", round, seq)
	}()

	t1 := time.Now()
	client, err := proxy.ServantByAddr(host + ":9001")
	if err != nil {
		report.incConnErrs()
		log.Printf("session{round^%d seq^%d} error: %v", round, seq, err)
		return
	}
	defer client.Recycle() // when err occurs, do we still need recycle?

	var enableLog = false
	if sampling(SampleRate) || Concurrency == 1 {
		enableLog = true
	}

	if enableLog {
		log.Printf("session{round^%d seq^%d} connected within %s",
			round, seq, time.Since(t1))
	}

	ctx := rpc.NewContext()
	ctx.Reason = "stress.go"
	ctx.Host = "stress.test.local"
	ctx.Ip = "127.0.0.1"
	for i := 0; i < LoopsPerSession; i++ {
		ctx.Rid = fmt.Sprintf("round:%d,seq:%d:i:%d,", round, seq, i+1)
		if Cmd&CallPing != 0 {
			var r string
			r, err = client.Ping(ctx)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d ping}: %v", round, seq, err)
				client.Close()
				return
			} else {
				report.incCallOk()
				if enableLog {
					log.Printf("session{round^%d seq^%d ping}: %v", round, seq, r)
				}
			}
		}

		if Cmd&CallLCache != 0 {
			key := fmt.Sprintf("lc_stress:%d", rand.Int())
			value := []byte("value of " + key)
			_, err = client.LcSet(ctx, key, value)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d lc_set} %v", round, seq, err)
				client.Close()
				return
			} else {
				report.incCallOk()
			}

			value, _, err = client.LcGet(ctx, key)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d lc_get} %v", round, seq, err)
				client.Close()
				return
			} else {
				report.incCallOk()
				if enableLog {
					log.Printf("session{round^%d seq^%d lcache}: %s => %s",
						round, seq, key, string(value))
				}
			}
		}

		if Cmd&CallIdGen != 0 {
			var r int64
			r, _, err = client.IdNext(ctx)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d idgen}: %v", round, seq, err)
				client.Close()
				return
			} else {
				if enableLog {
					log.Printf("session{round^%d seq^%d idgen}: %d",
						round, seq, r)
				}
				report.incCallOk()
			}
		}

		if Cmd&CallGame != 0 {
			client.GmLatency(ctx, 12, 4545)
			report.incCallOk()
			if enableLog {
				log.Printf("session{round^%d seq^%d GmLatency}",
					round, seq)
			}

			r, err := client.GmName3(ctx)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d GmName3}: %v", round, seq, err)
				client.Close()
				return
			} else {
				if enableLog {
					log.Printf("session{round^%d seq^%d GmName3}: %s",
						round, seq, r)
				}
				report.incCallOk()
			}

			lockKey := fmt.Sprintf("key:%d:%d", round, seq)
			client.GmLock(ctx, "stress.go", lockKey)
			report.incCallOk()
			if enableLog {
				log.Printf("session{round^%d seq^%d GmLock}: %s",
					round, seq, lockKey)
			}
			client.GmUnlock(ctx, "stress.go", lockKey)
			report.incCallOk()
			if enableLog {
				log.Printf("session{round^%d seq^%d GmUnlock}: %s",
					round, seq, lockKey)
			}

			shardId, err := client.GmRegister(ctx, "u")
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d GmRegister}: %v", round, seq, err)
				client.Close()
				return
			} else {
				if enableLog {
					log.Printf("session{round^%d seq^%d GmRegister}: %d",
						round, seq, shardId)
				}
				report.incCallOk()
			}
		}

		if Cmd&CallMysql != 0 {
			// with cache
			if true {
				r, err := client.MyQuery(ctx, "UserShard", "UserInfo", 1,
					"SELECT * FROM UserInfo WHERE uid=?",
					[]string{"1"}, "user:1")
				if err != nil {
					report.incCallErr()
					log.Printf("session{round^%d seq^%d mysql}: %v", round, seq, err)
					client.Close()
					return
				} else {
					if enableLog {
						log.Printf("session{round^%d seq^%d mysql}: %+v",
							round, seq, r)
					}
					report.incCallOk()
				}
			}

			// without cache
			if true {
				var rows *rpc.MysqlResult
				rows, err = client.MyQuery(ctx, "UserShard", "UserInfo", 1,
					"SELECT * FROM UserInfo WHERE uid=?",
					[]string{"1"}, "")
				if err != nil {
					report.incCallErr()
					log.Printf("session{round^%d seq^%d mysql}: %v", round, seq, err)
					client.Close()
					return
				} else {
					if enableLog {
						log.Printf("session{round^%d seq^%d mysql}: %+v",
							round, seq, rows.Rows)
					}
					report.incCallOk()
				}
			}

		}

		continue // TODO

		var (
			key     string
			value   []byte
			mcValue = rpc.NewTMemcacheData()
			result  []byte
		)
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		value = []byte("value of " + key)
		mcValue.Data = value

		if Cmd&CallMemcache != 0 {
			_, err = client.McSet(ctx, MC_POOL, key, mcValue, 36000)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d mc_set} %v", round, seq, err)
			} else {
				report.incCallOk()
			}
			_, miss, err := client.McGet(ctx, MC_POOL, key)
			if miss != nil || err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d mc_get} miss^%v, err^%v",
					round, seq, miss, err)
			} else {
				report.incCallOk()
			}
		}

		if Cmd&CallMongo != 0 {
			mgQuery, _ := bson.Marshal(bson.M{"snsid": "100003391571259"})
			mgFields, _ := bson.Marshal(bson.M{})
			result, _, err = client.MgFindOne(ctx, "default", "idmap",
				0, mgQuery, mgFields)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d mg_findOne} %v", round, seq, err)
			} else {
				report.incCallOk()

				if false {
					log.Println(result)
				}
			}
		}

	}

	if enableLog {
		log.Printf("session{round^%d seq^%d} done", round, seq)
	}

}
