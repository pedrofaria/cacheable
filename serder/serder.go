package serder

type Serder interface {
	Serialize(value any) ([]byte, error)
	Deserialize(data []byte, value any) error
}
