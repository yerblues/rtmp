package rtmp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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
	w   *bufio.Writer
	tmp [8]byte
	err error
}

func NewAMF3Encoder(w io.Writer) *AMF3Encoder {
	return NewAMF3EncoderSize(w, DefaultAMFBufSize)
}

func NewAMF3EncoderSize(w io.Writer, size int) *AMF3Encoder {
	return &AMF3Encoder{
		w: bufio.NewWriterSize(w, size),
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
		return e.encodeDouble(r.Float())
	case reflect.Bool:
		return e.encodeBool(r.Bool())
	case reflect.Map:
		return 0, nil
	default:
		return e.encodeUndefined()
	}
}

func (e *AMF3Encoder) encodeUndefined() (int, error) {
	return 1, e.w.WriteByte(byte(amf3DataTypeUndefined))
}

func (e *AMF3Encoder) encodeBool(data bool) (int, error) {
	if data {
		return 1, e.w.WriteByte(byte(amf3DataTypeTrue))
	}
	return 1, e.w.WriteByte(byte(amf3DataTypeFalse))
}

func (e *AMF3Encoder) encodeInteger(data uint32) (int, error) {
	if err := e.w.WriteByte(byte(amf3DataTypeInteger)); err != nil {
		return 0, err
	}
	var l int
	var err error
	switch {
	case data <= 0x7F:
		e.tmp[0] = byte(data)
		l, err = e.w.Write(e.tmp[:1])
	case 0x80 <= data && data <= 0x3FFF:
		e.tmp[0] = byte(data>>7 | 0x80)
		e.tmp[1] = byte(data & 0x7F)
		l, err = e.w.Write(e.tmp[:2])
	case 0x4000 <= data && data <= 0x1FFFFF:
		e.tmp[0] = byte(data>>14 | 0x80)
		e.tmp[1] = byte((data >> 7 & 0x7F) | 0x80)
		e.tmp[2] = byte(data & 0x7F)
		l, err = e.w.Write(e.tmp[:3])
	case 0x200000 <= data && data <= 0x1FFFFFFF:
		e.tmp[0] = byte((data >> 22 & 0x7F) | 0x80)
		e.tmp[1] = byte((data >> 15 & 0x7F) | 0x80)
		e.tmp[2] = byte((data >> 8 & 0x7F) | 0x80)
		e.tmp[3] = byte(data & 0xFF)
		l, err = e.w.Write(e.tmp[:4])
	default:
		return 0, ErrAMF3U29OverRange
	}
	return l + 1, err
}

func (e *AMF3Encoder) encodeDouble(data float64) (int, error) {
	if err := e.w.WriteByte(byte(amf3DataTypeDouble)); err != nil {
		return 0, err
	}
	bits := math.Float64bits(data)
	binary.BigEndian.PutUint64(e.tmp[:], bits)
	l, err := e.w.Write(e.tmp[:])
	return l + 1, err
}

func (e *AMF3Encoder) Flush() error {
	return e.w.Flush()
}

type AMF3Decoder struct {
	r   *bufio.Reader
	b   byte
	err error
}

func NewAMF3Decoder(r io.Reader) *AMF3Decoder {
	return NewAMF3DecoderSize(r, DefaultAMFBufSize)
}

func NewAMF3DecoderSize(r io.Reader, size int) *AMF3Decoder {
	return &AMF3Decoder{
		r: bufio.NewReaderSize(r, size),
	}
}

func (d *AMF3Decoder) Decode() (interface{}, error) {
	return d.decode()
}

func (d *AMF3Decoder) decode() (interface{}, error) {
	b, err := d.r.ReadByte()
	if err != nil {
		return amf3DataTypeUndefined, err
	}

	switch amf3DataType(b) {
	case amf3DataTypeUndefined:
	case amf3DataTypeNull:
	case amf3DataTypeFalse:
		return false, nil
	case amf3DataTypeTrue:
		return true, nil
	case amf3DataTypeInteger:
		return d.decodeInteger()
	case amf3DataTypeDouble:
		return d.decodeDouble()
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
		b, err := d.r.ReadByte()
		if err != nil {
			return 0, err
		}
		if i < 3 {
			data = (data << 7) + uint32(b&0x7F)
			if b&0x80 == 0 {
				break
			}
		} else {
			data = (data << 8) + uint32(b)
		}
	}
	return data, nil
}

func (d *AMF3Decoder) decodeDouble() (float64, error) {
	b, err := d.r.Peek(8)
	if err != nil {
		return 0, err
	}
	bits := binary.BigEndian.Uint64(b)
	if _, err := d.r.Discard(8); err != nil {
		return 0, err
	}
	return math.Float64frombits(bits), nil
}
