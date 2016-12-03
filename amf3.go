package rtmp

import (
	"fmt"
	"io"
	"reflect"
)

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

var (
	ErrAMF3U29OverRange = fmt.Errorf("U29 range error.")
)

type AMF3Encoder struct {
	io.Writer
	buf []byte
}

func NewAMF3Encoder(w io.Writer) *AMF3Encoder {
	return &AMF3Encoder{
		Writer: w,
		buf:    make([]byte, AMFBufSize),
	}
}

func (e *AMF3Encoder) Encode(data interface{}) (int, error) {
	return e.encode(data)
}

func (e *AMF3Encoder) encode(data interface{}) (int, error) {
	r := reflect.ValueOf(data)
	if !r.IsValid() {
		return 0, fmt.Errorf("invalid type.")
	}
	switch r.Kind() {
	case reflect.Uint32:
		return e.encodeInteger(uint32(r.Uint()))
	case reflect.Float64:
		return 0, nil
	case reflect.Bool:
		return e.encodeBool(r.Bool())
	case reflect.Map:
		return 0, nil
	default:
		return e.encodeUndefined()
	}
}

func (e *AMF3Encoder) encodeUndefined() (int, error) {
	e.appendMarker(amf3DataTypeUndefined)
	return e.Write(e.buf)
}

func (e *AMF3Encoder) encodeBool(data bool) (int, error) {
	if data {
		e.appendMarker(amf3DataTypeTrue)
	} else {
		e.appendMarker(amf3DataTypeFalse)
	}
	return e.Write(e.buf)
}

func (e *AMF3Encoder) encodeInteger(data uint32) (int, error) {
	e.appendMarker(amf3DataTypeInteger)
	switch {
	case data <= 0x7F:
		e.buf = append(e.buf,
			byte(data),
		)
	case 0x80 <= data && data <= 0x3FFF:
		e.buf = append(e.buf,
			byte(data>>7|0x80),
			byte(data&0x7F),
		)
	case 0x4000 <= data && data <= 0x1FFFFF:
		e.buf = append(e.buf,
			byte(data>>14|0x80),
			byte((data>>7&0x7F)|0x80),
			byte(data&0x7F),
		)
	case 0x200000 <= data && data <= 0x1FFFFFFF:
		e.buf = append(e.buf,
			byte((data>>22&0x7F)|0x80),
			byte((data>>15&0x7F)|0x80),
			byte((data>>8&0x7F)|0x80),
			byte(data&0xFF),
		)
	default:
		return 0, ErrAMF3U29OverRange
	}
	return e.Write(e.buf)
}

func (e *AMF3Encoder) appendMarker(marker amf3DataType) {
	e.buf = append(e.buf[:0], byte(marker))
}

type AMF3Decoder struct {
	io.Reader
	buf []byte
}

func NewAMF3Decoder(r io.Reader) *AMF3Decoder {
	return &AMF3Decoder{
		Reader: r,
		buf:    make([]byte, AMFBufSize),
	}
}

func (d *AMF3Decoder) Decode() (interface{}, error) {
	return d.decode()
}

func (d *AMF3Decoder) decode() (interface{}, error) {
	marker, err := d.readMarker()
	if err != nil {
		return amf3DataTypeUndefined, err
	}
	switch marker {
	case amf3DataTypeUndefined:
	case amf3DataTypeNull:
	case amf3DataTypeFalse:
		return false, nil
	case amf3DataTypeTrue:
		return true, nil
	case amf3DataTypeInteger:
		return d.decodeInteger()
	case amf3DataTypeDouble:
	case amf3DataTypeString:
	case amf3DataTypeXMLDoc:
	case amf3DataTypeDate:
	case amf3DataTypeArray:
	case amf3DataTypeObject:
	case amf3DataTypeXML:
	case amf3DataTypeByteArray:
	case amf3DataTypeVectorInt:
	case amf3DataTypeVectorUint:
	case amf3DataTypeVectorDouble:
	case amf3DataTypeVectorObject:
	case amf3DataTypeDictionary:
	default:
	}
	return nil, nil
}

func (d *AMF3Decoder) decodeInteger() (uint32, error) {
	var data uint32
	for i := 0; i < 4; i++ {
		if _, err := d.Read(d.buf[:1]); err != nil {
			return 0, err
		}
		if i < 3 {
			data = (data << 7) + uint32(d.buf[0]&0x7F)
			if d.buf[0]&0x80 == 0 {
				break
			}
		} else {
			data = (data << 8) + uint32(d.buf[0])
		}
	}
	return data, nil
}

func (d *AMF3Decoder) readMarker() (amf3DataType, error) {
	if _, err := d.Read(d.buf[:1]); err != nil {
		return amf3DataTypeUndefined, err
	}
	return amf3DataType(d.buf[0]), nil
}
