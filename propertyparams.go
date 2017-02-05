package util

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
)

type Params map[string]string

func ReadPropertyFile(fpath string) (o Params, err error) {
	f, err := os.Open(fpath)
	if nil != err {
		return
	}
	defer f.Close()

	return ReadPropertyParams(f)
}

// key / value separated by '='
func ReadPropertyParams(rdr io.Reader) (o Params, err error) {
	o = Params{}
	rb := bufio.NewReader(rdr)
	sep := []byte{'='}
	for {
		bbLine, _, err := rb.ReadLine()
		if io.EOF == err {
			err = nil
			break
		} else if nil != err {
			break
		}
		bbLine = bytes.TrimSpace(bbLine)
		if 0 == len(bbLine) {
			continue
		} else if '#' == bbLine[0] {
			// commented line
			continue
		}

		bbKV := bytes.SplitN(bbLine, sep, 2)
		if 2 != len(bbKV) {
			continue
		}

		v := string(bytes.TrimSpace(bbKV[1]))
		if '"' == v[0] && v[0] == v[len(v)-1] {
			v = v[1 : len(v)-1]
		}
		k := string(bytes.TrimSpace(bbKV[0]))
		o[strings.ToLower(k)] = v
	}

	return
}

func (this Params) Get(s, defval string) string {
	return this.getValue(s, defval).(string)
}

func (this Params) GetInt(s string, defval int) int {
	o := this.getValue(s, defval)
	v, ok := o.(int)
	if ok {
		return v
	}
	n, _ := strconv.ParseInt(o.(string), 10, 64)
	return int(n)
}

func (this Params) GetFloat64(s string, defval float64) float64 {
	o := this.getValue(s, defval)
	v, ok := o.(float64)
	if ok {
		return v
	}
	n, _ := strconv.ParseFloat(o.(string), 64)
	return n
}

func (this Params) getValue(s string, defaultValue interface{}) interface{} {
	v, has := this[strings.ToLower(s)]
	if !has {
		return defaultValue
	}

	return v
}
