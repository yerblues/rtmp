package rtmp

type Encoder interface {
	Encode(data interface{}) (int, error)
}

type Decoder interface {
	Decode() (interface{}, error)
}
