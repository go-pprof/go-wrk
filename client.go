package main

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"io"
	"sync"
)

func StartClient(u, h, m string, ch chan bool, dka bool, bc chan int64, rc chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := &http.Transport{DisableKeepAlives: dka}
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		panic(err)
	}
	sets := strings.Split(h, "\n")
	for i := range sets {
		split := strings.SplitN(sets[i], ":", 2)
		if len(split) == 2 {
			req.Header.Set(split[0], split[1])
		}
	}
	timer := NewTimer()
	for {
		timer.Reset()
		resp, err := tr.RoundTrip(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(ch) >= *totalCalls {
			break
		}
		_,err = io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			fmt.Println("error reading response:", err)
		}
		ch <- true
		rc <- resp.StatusCode
		bc <- timer.Duration()
		resp.Body.Close()
	}
}
