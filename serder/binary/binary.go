package binary

import (
	"bytes"
	"encoding/gob"
)

type binarySerder struct{}

func New() *binarySerder {
	return &binarySerder{}
}

func (b *binarySerder) Serialize(obj any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b *binarySerder) Deserialize(data []byte, obj any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	return dec.Decode(obj)
}
