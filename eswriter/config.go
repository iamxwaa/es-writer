package eswriter

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Mconfig struct {
	Ip        string `yaml:"ip"`
	Port      string `yaml:"port"`
	Open      bool   `yaml:"open"`
	Test      bool   `yaml:"switch-test"`
	Log       bool   `yaml:"switch-log"`
	StartTime int64
}

type Dconfig struct {
	Base     string `yaml:"base"`
	Pattern  string `yaml:"pattern"`
	Sotre    string `yaml:"store"`
	FileSize int    `yaml:"file-size"`
	WaitFile int    `yaml:"wait-file-count"`
	Interval int    `yaml:"interval"`
}

type Econfig struct {
	Batch   int      `yaml:"batch"`
	Nodes   []string `yaml:"nodes"`
	Master  string   `yaml:"master"`
	Port    string   `yaml:"port"`
	Delay   int64    `yaml:"delay"`
	Gzip    bool     `yaml:"switch-gzip"`
	TimeOut time.Duration
}

type Lconfig struct {
	Mc Mconfig `yaml:"metric-config"`
	Dc Dconfig `yaml:"device-file-config"`
	Ec Econfig `yaml:"elastic-search-config"`
}

var Lc Lconfig
var configPath string
var configModifyTime int64
var Config_reload_ch = make(chan bool, 1)
var Config_reload_mark bool

func WatchConfig(autoReload bool) {
	if autoReload {
		for true {
			fs, _ := os.Stat(configPath)
			if configModifyTime != fs.ModTime().Unix() {
				configModifyTime = fs.ModTime().Unix()
				Config_reload_mark = true
				Log("start reload config ...")
				<-Config_reload_ch
				Log("finish reload config ...")
				Config_reload_mark = false
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func LoadConfig() {
	Lc = Lconfig{}
	if "" == configPath {
		base := os.Args[0]
		i := strings.LastIndex(base, "/")
		if i == -1 {
			i = strings.LastIndex(base, "\\")
		}
		base = base[:i+1]
		Logf("base path : %s \r\n", base)
		configPath = base + "golc.yml"
	}
	fmt.Printf("config path : %s\n",configPath)
	f, e := os.Open(configPath)
	defer f.Close()
	fs, _ := f.Stat()
	configModifyTime = fs.ModTime().Unix()
	var c []byte = []byte("{}")
	if nil == e {
		Logf("config path : %s \r\n", configPath)
		c, _ = ioutil.ReadAll(f)
	}

	yaml.Unmarshal(c, &Lc)

	setDc(&Lc.Dc)
	setMc(&Lc.Mc)
	setEc(&Lc.Ec)

	LogC(Lc)
}

func setDc(dc *Dconfig) {
	if "" == dc.Pattern {
		dc.Pattern = ".*(Z_[0-9]{2}).*"
	}
	if "" == Lc.Dc.Sotre {
		dc.Sotre = strings.Join([]string{filepath.Dir(os.Args[0]), string(filepath.Separator), "wal", string(filepath.Separator)}, "")
		os.Mkdir(dc.Sotre, 777)
	}
	if 0 == dc.FileSize {
		dc.FileSize = 1024 * 1024 * 10
	}
	if 0 == dc.WaitFile {
		dc.WaitFile = 5
	}
	if 2 > dc.WaitFile {
		dc.WaitFile = 2
	}
	if 0 == dc.Interval {
		dc.Interval = 10
	}
}

func setEc(ec *Econfig) {
	if 0 == ec.Batch {
		ec.Batch = 5000
	}
	if 10000 < ec.Batch {
		ec.Batch = 10000
	}
	if 0 == ec.Delay {
		ec.Delay = int64(1000)
	}
	if 0 == ec.TimeOut {
		ec.TimeOut = 10 * time.Second
	} else {
		ec.TimeOut = ec.TimeOut * time.Second
	}
	if "" != ec.Master {
		if "" == ec.Port {
			ec.Port = ec.Master[strings.LastIndex(ec.Master, ":")+1:]
		}
		ReloadNodes()
	}
}

func setMc(mc *Mconfig) {
	if "" == mc.Ip {
		mc.Ip = "0.0.0.0"
	}
	if "" == mc.Port {
		mc.Port = "10000"
	}
	mc.StartTime = time.Now().UnixNano()
}

func ReloadNodes() {
	client := http.Client{}
	client.Timeout = Lc.Ec.TimeOut
	req := &http.Request{}
	req.Method = "GET"
	req.URL, _ = url.Parse(Lc.Ec.Master + "/_cat/nodes")
	req.Header = map[string][]string{"Accept-Encoding": []string{""}, "User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"}}
	resp, err := client.Do(req)
	//	resp, err := client.Get()
	if nil != err {
		Log(err)
	}
	if nil == resp {
		return
	}
	defer resp.Body.Close()
	c := make([]byte, resp.ContentLength)
	resp.Body.Read(c)
	d := bytes.Split(c, []byte("\n"))
	nlen := len(d)
	if "" == string(d[nlen-1]) {
		nlen -= 1
	}
	var nodes = make([]string, nlen)
	for i := 0; i < len(nodes); i++ {
		line := string(d[i])
		if "" == line {
			break
		}
		ip := line[:strings.Index(line, " ")]
		b := bytes.NewBuffer(make([]byte, 0, 20+len(line)))
		b.WriteString("http://")
		b.WriteString(string(ip))
		b.WriteString(":")
		b.WriteString(Lc.Ec.Port)
		b.WriteString("/_bulk")
		nodes[i] = b.String()
	}
	Lc.Ec.Nodes = nodes
}
