package resp

import (
	"bufio"
	"io"
)

func ParseResponse(rd io.Reader) (interface{}, error) {
	brd := bufio.NewReader(rd)
	prefix, err := brd.ReadByte()
	if err != nil {
		return nil, err
	}
	switch prefix {
	case _typeSimplePrefix:
		return decodeSimpleString(brd)
	case _typeIntegerPrefix:
		return decodeInteger(brd)
	case _typeBulkPrefix:
		return decodeBulk(brd)
	case _typeArrayPrefix:
		return decodeArray(brd)
	case _typeErrorPrefix:
		return decodeError(brd)
	default:
		return nil, errIllegalProto
	}
	return nil, nil
}
