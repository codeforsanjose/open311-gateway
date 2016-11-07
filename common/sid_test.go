package common_test

import (
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common"
)

type sidtest struct {
	m map[int64]bool
	sync.Mutex
}

var _ = Describe("Common Suite", func() {
	Describe("RequestID", func() {
		It("run several goroutines", func() {
			stest := sidtest{
				m: make(map[int64]bool),
			}
			rand.Seed(102384)
			var wg sync.WaitGroup
			wg.Add(3)

			f := func() {
				var ri int64
				for i := 0; i < 30; i++ {
					stest.Lock()
					ri = RequestID()
					stest.m[ri] = true
					stest.Unlock()
					time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				}
				wg.Done()
			}

			go f()
			go f()
			go f()

			wg.Wait()

			var i int64
			for i = 101; i < 190; i++ {
				// fmt.Printf("m[%v] %v\n", i, stest.m[i])
				Ω(stest.m[i]).Should(BeTrue())
			}
		})

		Describe("RPCID", func() {
			It("run several goroutines", func() {
				stest := sidtest{
					m: make(map[int64]bool),
				}
				rand.Seed(102384)
				var wg sync.WaitGroup
				wg.Add(3)

				f := func() {
					var ri int64
					for i := 0; i < 30; i++ {
						stest.Lock()
						ri = RPCID()
						stest.m[ri] = true
						stest.Unlock()
						time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
					}
					wg.Done()
				}

				go f()
				go f()
				go f()

				wg.Wait()

				var i int64
				for i = 2; i < 90; i++ {
					// fmt.Printf("m[%v] %v\n", i, stest.m[i])
					Ω(stest.m[i]).Should(BeTrue())
				}
			})
		})

	})
})
