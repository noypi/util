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
	var namespaces []string
	var ns string

	for {
		bbLine, _, err := rb.ReadLine()
		if io.EOF == err {
			err = nil
			break
		} else if nil != err {
			break
		}
		var bIsEnvVar bool
		bbLine = bytes.TrimSpace(bbLine)
		if 0 == len(bbLine) {
			continue
		} else if '#' == bbLine[0] {
			// commented line
			continue
		} else if '$' == bbLine[0] {
			bIsEnvVar = true
		} else if bytes.HasPrefix(bbLine, []byte(":copy")) {
			if 0 == len(ns) {
				// invalid, should be in a namespace
				// ignore
				continue
			}

			srcNs := string(bytes.TrimSpace(bytes.TrimPrefix(bbLine, []byte(":copy"))))
			if ns == srcNs {
				// same origin, ignore
				continue
			}

			o.For(srcNs, func(k, v string) {
				k = strings.TrimPrefix(k, srcNs)
				o[ns[0:len(ns)-1]+k] = v
			})

			continue

		} else if bytes.HasPrefix(bbLine, []byte(":source")) {
			if 0 < len(namespaces) {
				// ignore, source cannot be used inside of a namespace
				continue
			}
			bbSource := bytes.TrimSpace(bytes.TrimPrefix(bbLine, []byte(":source")))
			if 0 == len(bbSource) {
				// invalid ignore
				continue
			}
			if o2, err := ReadPropertyFile(string(bbSource)); nil != err {
				return o, err
			} else {
				o.Merge(o2)
			}
			continue

		} else if nNs := countPrefix(bbLine, '['); 0 < nNs {
			// namespace
			if !bytes.HasSuffix(bbLine, bytes.Repeat([]byte{']'}, nNs)) {
				// invalid, ignore
				continue
			}

			if len(namespaces) >= nNs {
				// truncate, we are in new namespace
				namespaces = namespaces[0 : nNs-1]
			}

			if len(namespaces)+1 != nNs {
				// invalid, ignore
				continue
			}

			currNs := bbLine[nNs:]
			currNs = currNs[0 : len(currNs)-nNs]
			namespaces = append(namespaces, string(currNs))
			ns = strings.Join(namespaces, ".") + "."
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
		v = os.Expand(v, o.Env)

		if 0 < len(ns) {
			k = ns + k
		}

		if bIsEnvVar {
			o[k] = v
		} else {
			o[strings.ToLower(k)] = v
		}
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

func (this Params) Merge(a Params) {
	for k, v := range a {
		if '$' == k[0] {
			continue
		}
		this[k] = os.Expand(v, this.Env)
	}
}

func (this Params) For(prefix string, cb func(k, v string)) {
	for k, v := range this {
		if strings.HasPrefix(k, prefix) {
			cb(k, v)
		}
	}
}

func (this Params) Env(k string) string {
	s, _ := this["$"+k]
	return s
}

func countPrefix(s []byte, c byte) int {
	for i, c0 := range s {
		if c != c0 {
			return i
		}
	}
	return 0
}
