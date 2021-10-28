package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type decodingFunction func(decoder *decoder) (interface{}, error)

type decoder struct {
	reader *bufio.Reader
	delim  byte
}

func newDecoder() decoder {
	return decoder{}
}

//Decode: decode benconded data
func Decode(reader io.Reader) (map[interface{}]interface{}, error) {
	decoder := newDecoder()
	decoder.delim = 'e'
	decoder.reader = bufio.NewReader(reader)
	first, err := decoder.reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if first != 'd' {
		return nil, errors.New("first value in bencoded data must be dictionary")
	}
	decodedDict, err := decodeDict(&decoder)
	if err != nil {
		return nil, err
	}
	return decodedDict.(map[interface{}]interface{}), nil
}

func (decoder *decoder) fetchFunctionForNextType() (decodingFunction, error) {
	byte, err := decoder.reader.ReadByte()
	if err != nil {
		return nil, err
	}
	switch byte {
	case 'd':
		return decodeDict, nil
		break
	case 'l':
		return readList, nil
		break
	case 'i':
		return readInt, nil
		break
	default:
		//can only be string
		decoder.reader.UnreadByte()
		return readString, nil
	}
	return nil, errors.New("cannot determine function for type")
}

func decodeDict(decoder *decoder) (interface{}, error) {
	values := make(map[interface{}]interface{})
	for {
		byte, err := decoder.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if byte == 'e' {
			return values, nil
		}
		decoder.reader.UnreadByte()
		key, err := readString(decoder)
		if err != nil {
			return nil, fmt.Errorf("cannot read dict key: %w", err)
		}
		decodingFunc, err := decoder.fetchFunctionForNextType()
		if err != nil {
			return nil, err
		}
		val, err := decodingFunc(decoder)
		if err != nil {
			return nil, err
		}
		values[key] = val
	}
}

func readString(decoder *decoder) (interface{}, error) {
	decoder.delim = ':'
	strSizeI, err := readInt(decoder)
	decoder.delim = 'e'
	if err != nil {
		return "", err
	}
	strSize := strSizeI.(int64)
	if strSize < 0 {
		return "", errors.New(fmt.Sprintf("invalid byte string size %d", strSize))
	} else if strSize == 0 {
		return "", nil
	}
	//add space for :
	buf := make([]byte, strSize)
	n, err := decoder.reader.Read(buf)
	if err != nil {
		return "", err
	}
	if int64(n) != (strSize) {
		return "", errors.New("invalid read of string, total size and read size mismatch")
	}
	return string(buf[0:strSize]), nil
}

func readInt(decoder *decoder) (interface{}, error) {
	intData, err := decoder.reader.ReadBytes(decoder.delim)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(intData[0:len(intData)-1]), 10, 64) //ASCII byte numbers
	if err == nil {
		return n, nil
	}
	if err != nil && errors.Is(err, strconv.ErrRange) {
		un, err := strconv.ParseUint(string(intData[0:len(intData)-1]), 10, 64)
		if err == nil {
			return un, nil
		}
	}
	return 0, err
}

func readList(decoder *decoder) (interface{}, error) {
	list := make([]interface{}, 0)
	for {
		byte, err := decoder.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if byte == decoder.delim {
			break
		}
		decoder.reader.UnreadByte()
		decodingFunc, err := decoder.fetchFunctionForNextType()
		if err != nil {
			return nil, err
		}
		res, err := decodingFunc(decoder)
		if err != nil {
			return nil, err
		}
		list = append(list, res)
	}
	return list, nil
}
