package msgpack

import "github.com/vmihailenco/msgpack/v5"

type msgpackSerder struct{}

func New() *msgpackSerder {
	return &msgpackSerder{}
}

func (b *msgpackSerder) Serialize(obj any) ([]byte, error) {
	return msgpack.Marshal(obj)
}

func (b *msgpackSerder) Deserialize(data []byte, obj any) error {
	return msgpack.Unmarshal(data, obj)
}
