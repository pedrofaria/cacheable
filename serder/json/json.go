package json

import "encoding/json"

type jsonSerder struct{}

func New() *jsonSerder {
	return &jsonSerder{}
}

func (j *jsonSerder) Serialize(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func (j *jsonSerder) Deserialize(data []byte, obj any) error {
	return json.Unmarshal(data, obj)
}
