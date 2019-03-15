package eswriter

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

func Log(v ...interface{}) {
	if Lc.Mc.Log {
		log.Println(v)
	}
}

func Logf(format string, v ...interface{}) {
	if Lc.Mc.Log {
		log.Printf(format, v...)
	}
}

func LogC(lc Lconfig) {
	fmt.Fprintf(os.Stderr, "************************************************\n\n")
	rv := reflect.ValueOf(lc)
	for i := 0; rv.IsValid() && i < rv.NumField(); i++ {
		v := rv.Field(i)
		fname := rv.Type().Field(i).Name
		for j := 0; j < v.NumField(); j++ {
			t := v.Type().Field(j)
			f := v.Field(j)
			fmt.Fprintf(os.Stderr, "  Lc.%s.%s = %v\n", fname, t.Name, f)
		}
	}
	fmt.Fprintf(os.Stderr, "\n************************************************\n")
}
