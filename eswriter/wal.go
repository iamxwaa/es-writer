package eswriter

import (
	"bufio"
	"bytes"
	"github.com/pquerna/ffjson/ffjson"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	JDATA  = "jdata"
	JDATAF = "Jdata"
	JTAG   = "json"

	WAL_TEMP = ".temp"
	WAL_LOG  = ".log"
)

var (
	decodes = []string{"CreateTime", "ActiveTime", "DisableTime", "path", "time", "uri", "referer", "Time", "WinF"}

	Z30      = IndexInfo{Entry: &Net{}, Index: "test-a-net-", TimeKey: "Time", Type: "Z_30"}
	Z31      = IndexInfo{Entry: &StartShutdown{}, Index: "test-a-startshutdown-", TimeKey: "Time", Type: "Z_31"}
	Z32      = IndexInfo{Entry: &Process{}, Index: "test-a-pro-", TimeKey: "CreateTime", Type: "Z_32"}
	Z43      = IndexInfo{Entry: &Copy{}, Index: "test-a-copy-", TimeKey: "Time", Type: "Z_43"}
	IndexMap = map[string]IndexInfo{"Z_30": Z30, "Z_31": Z31, "Z_32": Z32, "Z_43": Z43}

	indexBuffer = bytes.NewBuffer(make([]byte, 0, 128))
	headBuffer  = bytes.NewBuffer(make([]byte, 0, 512))
	empty       = make([]byte, 0)
	kv          = Kv{}
	logWriter   LogWriter
)

type Kv struct {
	Field string
	Value string
}

type LogWriter struct {
	writer  *os.File
	wbuffer *bufio.Writer
	size    int
	path    string
}

func (logWriter *LogWriter) newWriter() {
	if nil == logWriter.writer {
		logWriter.newPath()
		logWriter.writer, _ = os.Create(logWriter.path)
		logWriter.wbuffer = bufio.NewWriter(logWriter.writer)
		logWriter.size = 0
	}
}

func (logWriter *LogWriter) Write(b []byte) {
	logWriter.newWriter()
	logWriter.wbuffer.Write(b)
	logWriter.size += len(b)
}

func (logWriter *LogWriter) WriteString(s string) {
	logWriter.newWriter()
	logWriter.wbuffer.WriteString(s)
	logWriter.size += len(s)
}

func (logWriter *LogWriter) FlushWrite(s string) {
	logWriter.wbuffer.WriteString(s)
	logWriter.size += len(s)
	if Lc.Dc.FileSize <= logWriter.size {
		logWriter.wbuffer.Flush()
		logWriter.writer.Close()
		newName := strings.Join([]string{logWriter.path[:len(logWriter.path)-5], WAL_LOG}, "")
		os.Rename(logWriter.path, newName)
		logWriter.newPath()
		logWriter.writer, _ = os.Create(logWriter.path)
		logWriter.wbuffer = bufio.NewWriter(logWriter.writer)
		logWriter.size = 0
	}
}

func (logWriter *LogWriter) newPath() {
	logWriter.path = strings.Join([]string{Lc.Dc.Sotre, time.Now().UTC().Format("20060102150405"), WAL_TEMP}, "")
}

func WriteBeforeSend(c []byte, ftype string) {
	indexInfo := IndexMap[ftype]
	if nil != indexInfo.Entry {
		ffjson.UnmarshalFast(c, indexInfo.Entry)
		header := getCommonHeader(indexInfo)
		writeToLog(indexInfo, header)
	}
}

func writeToLog(indexInfo IndexInfo, head []byte) {
	v2 := reflect.ValueOf(indexInfo.Entry).Elem()
	jdata := v2.FieldByName(JDATAF)
	var index []byte
	for i := 0; i < jdata.Len(); i++ {
		v3 := jdata.Index(i)
		index = getIndex(v3, indexInfo)
		if 0 == len(index) {
			continue
		}
		logWriter.Write(index)
		logWriter.WriteString("\r\n")
		logWriter.WriteString("{")
		logWriter.Write(head)
		for i := 0; i < v3.NumField(); i++ {
			kv.toKv(v3, i)
			logWriter.WriteString("\"")
			logWriter.WriteString(kv.Field)
			logWriter.WriteString("\":\"")
			logWriter.WriteString(kv.Value)
			if i == (v3.NumField() - 1) {
				logWriter.WriteString("\"")
			} else {
				logWriter.WriteString("\",")
			}
		}
		logWriter.FlushWrite("}\r\n")
	}
}

func getCommonHeader(indexInfo IndexInfo) []byte {
	headBuffer.Reset()
	vv := reflect.ValueOf(indexInfo.Entry).Elem()
	for i := 0; i < vv.NumField(); i++ {
		kv.toKv(vv, i)
		if JDATA == kv.Field {
			continue
		}
		headBuffer.WriteString("\"")
		headBuffer.WriteString(kv.Field)
		headBuffer.WriteString("\":\"")
		headBuffer.WriteString(kv.Value)
		headBuffer.WriteString("\",")
	}
	return headBuffer.Bytes()
}

func getIndex(v reflect.Value, indexInfo IndexInfo) []byte {
	indexBuffer.Reset()
	time := v.FieldByName(indexInfo.TimeKey).String()
	if len(time) < 10 {
		return empty
	}
	time = strings.Replace(time[0:10], "-", ".", 2)
	indexBuffer.WriteString("{\"index\":{\"_index\":\"")
	indexBuffer.WriteString(indexInfo.Index)
	indexBuffer.WriteString(time)
	indexBuffer.WriteString("\",\"_type\":\"logs\"}}")
	return indexBuffer.Bytes()
}

func (kv *Kv) toKv(v reflect.Value, i int) {
	f := v.Type().Field(i)
	kv.Field = f.Tag.Get(JTAG)
	if "" == kv.Field {
		kv.Field = f.Name
	}
	kv.Value = v.Field(i).String()
	for i := 0; i < len(decodes); i++ {
		if kv.Field == decodes[i] {
			kv.Value, _ = url.PathUnescape(kv.Value)
			break
		}
	}
}
