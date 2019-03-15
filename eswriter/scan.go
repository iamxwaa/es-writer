package eswriter

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ZIP    = ".zip"
	REPORT = "report.dat"
)

var (
	r             *regexp.Regexp
	contentBuffer = bytes.NewBuffer(make([]byte, 0, 1024*1024*2))
)

func Init() {
	r, _ = regexp.Compile(Lc.Dc.Pattern)
}

func Scan(base string) {
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if Config_reload_mark {
			return nil
		}
		content, ftype := getContent(path, info)
		if 0 != len(content) {
			WriteBeforeSend(content, ftype)
		}
		return nil
	})
}

func getContent(path string, info os.FileInfo) ([]byte, string) {
	contentBuffer.Reset()
	var ftype []string
	if !info.IsDir() && strings.HasSuffix(info.Name(), ZIP) {
		ftype = r.FindStringSubmatch(info.Name())
		if len(ftype) > 1 {
			reader, _ := zip.OpenReader(path)
			defer reader.Close()
			for _, file := range reader.File {
				if REPORT == file.Name {
					r, _ := file.Open()
					defer r.Close()
					io.Copy(contentBuffer, r)
				}
			}
			return contentBuffer.Bytes(), ftype[1]
		}
	}
	return make([]byte, 0), ""
}
