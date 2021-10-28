package bencode

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type encoder struct {
	buffer *bytes.Buffer
}

func newEncoder() encoder {
	return encoder{
		buffer: new(bytes.Buffer),
	}
}

//Encode: encode bencoded data
func Encode(value interface{}) []byte {
	encoder := newEncoder()
	encoder.parseType(value)
	return encoder.buffer.Bytes()
}

func (encoder *encoder) parseType(value interface{}) {
	switch value.(type) {
	case string:
		encoder.encodeByteString(value.(string))
		break
	case map[string]interface{}:
		encoder.encodeDict(value.(map[string]interface{}))
		break
	case []interface{}:
		encoder.encodeList(value.([]interface{}))
		break
	case int8, int16, int32, int, int64:
		encoder.encodeInt(reflect.ValueOf(value).Int())
		break
	case uint8, uint16, uint32, uint64:
		encoder.encodeUint(reflect.ValueOf(value).Uint())
		break
	}
}

func (encoder *encoder) encodeDict(value map[string]interface{}) {
	//sort the keys into a list because iteration order of maps is random in go
	list := make(sort.StringSlice, len(value))
	i := 0
	for k := range value {
		list[i] = k
		i++
	}
	list.Sort()
	encoder.buffer.WriteByte('d')
	for i := 0; i < len(list); i++ {
		encoder.encodeByteString(list[i])
		encoder.parseType(value[list[i]])
	}
	encoder.buffer.WriteByte('e')
}

func (encoder *encoder) encodeList(list []interface{}) {
	encoder.buffer.WriteByte('l')
	for i := 0; i < len(list); i++ {
		encoder.parseType(list[i])
	}
	encoder.buffer.WriteByte('e')
}

func (encoder *encoder) encodeByteString(value string) {
	encoder.buffer.WriteString(fmt.Sprintf("%d", len(value)))
	encoder.buffer.WriteByte(':')
	encoder.buffer.WriteString(value)
}

func (encoder *encoder) encodeInt(value int64) {
	encoder.buffer.WriteByte('i')
	encoder.buffer.WriteString(fmt.Sprintf("%d", value))
	encoder.buffer.WriteByte('e')
}

func (encoder *encoder) encodeUint(value uint64) {
	encoder.buffer.WriteByte('i')
	encoder.buffer.WriteString(fmt.Sprintf("%d", value))
	encoder.buffer.WriteByte('e')
}
