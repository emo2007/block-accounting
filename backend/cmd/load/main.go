package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/pkg/bip39"
)

func main() {
	sAt := time.Now()

	wg := sync.WaitGroup{}

	totalch := make(chan int, 5)

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(j int) {
			defer wg.Done()

			var reqc int

			for {
				e, err := bip39.NewEntropy(256)
				if err != nil {
					log.Println("ERROR: ", err)
					break
				}

				m, err := bip39.NewMnemonic(e)
				if err != nil {
					log.Println("ERROR: ", err)
					break
				}

				req, err := json.Marshal(&domain.JoinRequest{
					Mnemonic: m,
				})
				if err != nil {
					log.Println("ERROR: ", err)
					break
				}

				_, err = http.Post("http://localhost:8080/join", "application/json", bytes.NewBuffer(req))
				if err != nil {
					log.Println("ERROR: ", err)
					break
				}

				reqc++

				log.Println("req ", j)
			}

			totalch <- reqc
		}(i)
	}

	var reqtotoal int
	mu := sync.Mutex{}

	go func() {
		for c := range totalch {
			mu.Lock()
			reqtotoal += c
			mu.Unlock()
		}
	}()

	wg.Wait()

	eAt := time.Now()

	rps := float64(reqtotoal) / eAt.Sub(sAt).Seconds()

	log.Println("STARTED_AT: ", sAt, " END_AT: ", eAt)

	log.Println("REQ TOTAL: ", reqtotoal, " RPS:", rps, " SECONDS: ", eAt.Sub(sAt).Seconds())
}
