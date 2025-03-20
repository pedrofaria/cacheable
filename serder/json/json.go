package json

import "encoding/json"

type JsonSerde struct{}

func NewJsonSerde() *JsonSerde {
	return &JsonSerde{}
}

func (j *JsonSerde) Serialize(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func (j *JsonSerde) Deserialize(data []byte, obj any) error {
	return json.Unmarshal(data, obj)
}
