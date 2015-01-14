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
	}()

	t1 := time.Now()
	client, err := proxy.Servant(host + ":9001")
	if err != nil {
		report.incConnErrs()
		log.Printf("session{round^%d seq^%d} error: %v", round, seq, err)
		return
	}
	defer client.Recycle() // when err occurs, do we still need recycle?

	var enableLog = false
	if sampling(SampleRate) {
		enableLog = true
	}

	if enableLog {
		log.Printf("session{round^%d seq^%d} connected within %s",
			round, seq, time.Since(t1))
	}

	var (
		key     string
		value   []byte
		mcValue = rpc.NewTMemcacheData()
		result  []byte
	)
	mgQuery, _ := bson.Marshal(bson.M{"snsid": "100003391571259"})
	mgFields, _ := bson.Marshal(bson.M{})
	for i := 0; i < LoopsPerSession; i++ {
		if Cmd&CallPing != 0 {
			_, err = client.Ping(ctx)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d ping} %v", round, seq, err)
			} else {
				report.incCallOk()
			}
		}

		if Cmd&CallIdGen != 0 {
			_, _, err = client.IdNext(ctx)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d idgen} %v", round, seq, err)
			} else {
				report.incCallOk()
			}
		}

		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		value = []byte("value of " + key)
		mcValue.Data = value

		if Cmd&CallLCache != 0 {
			_, err = client.LcSet(ctx, key, value)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d lc_set} %v", round, seq, err)
			} else {
				report.incCallOk()
			}
			_, _, err = client.LcGet(ctx, key)
			if err != nil {
				report.incCallErr()
				log.Printf("session{round^%d seq^%d lc_get} %v", round, seq, err)
			} else {
				report.incCallOk()
			}
		}

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
