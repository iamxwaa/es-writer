package eswriter

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	client         = &http.Client{}
	nodesIndex int = 0
	head           = map[string][]string{
		"Content-Type": []string{"application/json;charset=UTF-8"},
		"User-Agent":   []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"}}
	post_ch = make(chan bool, 2)
	post_wg sync.WaitGroup
)

func Watch(base string) {
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if Config_reload_mark {
			return nil
		}
		if nil != info {
			if strings.HasSuffix(info.Name(), WAL_LOG) {
				post_wg.Add(1)
				go sendToEs(path, false)
			} else if strings.HasSuffix(info.Name(), WAL_TEMP) {
				if info.ModTime().UnixNano() < Lc.Mc.StartTime {
					post_wg.Add(1)
					go sendToEs(path, true)
				}
			}
		}
		return nil
	})
	post_wg.Wait()
}

func sendToEs(path string, temp bool) {
	post_ch <- true
	if temp {
		Logf("send old temp to es: %s\n", path)
	} else {
		Logf("send to es: %s\n", path)
	}
	if Lc.Mc.Test {
		ss := time.Duration(Lc.Dc.FileSize/(1024*1024*2) - rand.Intn(2))
		time.Sleep(ss * time.Second)
	} else {
		file, _ := os.Open(path)
		defer file.Close()
		req := getReq(file)
		resp, _ := client.Do(req)
		if nil != resp {
			defer resp.Body.Close()
		}
	}
	Logf("remove: %s\n", path)
	os.Remove(path)
	post_wg.Done()
	<-post_ch
}

func getReq(file *os.File) *http.Request {
	node := Lc.Ec.Nodes[chooseNodes()]
	req, _ := http.NewRequest("POST", node, file)
	req.Header = head
	fstat, _ := file.Stat()
	req.Header["Content-Length"] = []string{strconv.FormatInt(fstat.Size(), 10)}
	return req
}

func chooseNodes() int {
	if 0 == len(Lc.Ec.Nodes) {
		ReloadNodes()
		nodesIndex = len(Lc.Ec.Nodes) - 1
		return nodesIndex
	}
	nodesIndex--
	if nodesIndex < 0 {
		nodesIndex = len(Lc.Ec.Nodes) - 1
	}
	return nodesIndex
}
