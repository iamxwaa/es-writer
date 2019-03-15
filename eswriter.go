package main

import (
	"io/ioutil"
	. "es-writer/eswriter"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"
)

func main() {
loop:
	LoadConfig()
	Init()
	go WatchConfig(true)
	if Lc.Mc.Open {
		go func() {
			u := strings.Join([]string{Lc.Mc.Ip, ":", Lc.Mc.Port}, "")
			Logf("pprof listen at %s\r\n", u)
			Log(http.ListenAndServe(u, nil))
		}()
	}

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for !Config_reload_mark {
			fs, _ := ioutil.ReadDir(Lc.Dc.Sotre)
			if len(fs) < Lc.Dc.WaitFile {
				Scan(Lc.Dc.Base)
			}
			time.Sleep(time.Duration(Lc.Dc.Interval) * time.Second)
		}
		Log("stop scan file")
		wait.Done()
	}()
	wait.Add(1)
	go func() {
		for !Config_reload_mark {
			Watch(Lc.Dc.Sotre)
			time.Sleep(10 * time.Second)
		}
		Log("stop watch log file")
		wait.Done()
	}()
	wait.Wait()
	Log("prepare restart application")
	Config_reload_ch <- false
	goto loop
}
