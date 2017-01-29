package util

import (
	"bytes"
	"encoding/gob"
)

func SerializeGob(a interface{}) ([]byte, error) {
	buf := bytes.NewBufferString("")
	err := gob.NewEncoder(buf).Encode(a)
	return buf.Bytes(), err
}

func DeserializeGob(a interface{}, bb []byte) error {
	buf := bytes.NewBuffer(bb)
	return gob.NewDecoder(buf).Decode(a)
}
