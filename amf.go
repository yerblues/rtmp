package rtmp

const DefaultAMFBufSize = 1024

type AMFEncoder interface {
	Encode(data interface{}) (int, error)
}

type AMFDecoder interface {
	Decode() (interface{}, error)
}
