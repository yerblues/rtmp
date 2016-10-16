package rtmp

import "io"

// SEE: http://wwwimages.adobe.com/www.adobe.com/content/dam/Adobe/en/devnet/amf/pdf/amf-file-format-spec.pdf

type amf3DataType uint8

const (
	amf3DataTypeUndefined    amf3DataType = 0x00
	amf3DataTypeNull         amf3DataType = 0x01
	amf3DataTypeFalse        amf3DataType = 0x02
	amf3DataTypeTrue         amf3DataType = 0x03
	amf3DataTypeInteger      amf3DataType = 0x04
	amf3DataTypeDouble       amf3DataType = 0x05
	amf3DataTypeString       amf3DataType = 0x06
	amf3DataTypeXMLDoc       amf3DataType = 0x07
	amf3DataTypeDate         amf3DataType = 0x08
	amf3DataTypeArray        amf3DataType = 0x09
	amf3DataTypeObject       amf3DataType = 0x0a
	amf3DataTypeXML          amf3DataType = 0x0b
	amf3DataTypeByteArray    amf3DataType = 0x0c
	amf3DataTypeVectorInt    amf3DataType = 0x0d
	amf3DataTypeVectorUint   amf3DataType = 0x0e
	amf3DataTypeVectorDouble amf3DataType = 0x0f
	amf3DataTypeVectorObject amf3DataType = 0x10
	amf3DataTypeDictionary   amf3DataType = 0x11
)

type AMF3Encoder struct {
	io.Writer
	buf []byte
}

func NewAMF3Encoder(w io.Writer) *AMF3Encoder {
	return &AMF3Encoder{Writer: w}
}

func (e *AMF3Encoder) Encode(data interface{}) (int, error) {
	return e.encode(data)
}

func (e *AMF3Encoder) encode(data interface{}) (int, error) {
	switch d := data.(type) {
	case uint64:
		return 0, nil
	case float64:
		return 0, nil
	case bool:
		return e.encodeBool(d)
	case []byte:
		return 0, nil
	default:
		return e.encodeUndefined()
	}
}

func (e *AMF3Encoder) encodeUndefined() (int, error) {
	return e.Write([]byte{byte(amf3DataTypeUndefined)})
}

func (e *AMF3Encoder) encodeBool(data bool) (int, error) {
	if data {
		return e.Write([]byte{byte(amf3DataTypeTrue)})
	}
	return e.Write([]byte{byte(amf3DataTypeFalse)})
}
