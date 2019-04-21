package resp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type SimpleString string
type Error error
type Integer int
type Bulk []byte
type Array []interface{}

func encodeSimpleString(w io.Writer, val SimpleString) error {
	_, err := w.Write([]byte{_typeSimplePrefix})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(val))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(_lineDelimiter))
	return err
}

func decodeSimpleString(reader *bufio.Reader) (SimpleString, error) {
	body, err := reader.ReadString('\n')
	if err != nil {
		return SimpleString(""), err
	}
	body = strings.TrimPrefix(body, string(_typeSimplePrefix))
	body = strings.TrimSuffix(body, _lineDelimiter)
	return SimpleString(body), nil
}

func encodeError(w io.Writer, val Error) error {
	_, err := w.Write([]byte{_typeErrorPrefix})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(val.Error()))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(_lineDelimiter))
	return err
}

func decodeError(reader *bufio.Reader) (Error, error) {
	body, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	body = strings.TrimPrefix(body, string(_typeSimplePrefix))
	body = strings.TrimSuffix(body, _lineDelimiter)
	return errors.New(body), nil
}

func encodeInteger(w io.Writer, val Integer) error {
	_, err := w.Write([]byte{_typeIntegerPrefix})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(strconv.Itoa(int(val))))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(_lineDelimiter))
	return err
}

func decodeInteger(reader *bufio.Reader) (Integer, error) {
	body, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	body = strings.TrimPrefix(body, string(_typeIntegerPrefix))
	body = strings.TrimSuffix(body, _lineDelimiter)
	i, err := strconv.Atoi(body)
	return Integer(i), err
}

func encodeBulk(w io.Writer, val Bulk) error {
	_, err := w.Write([]byte{_typeBulkPrefix})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(strconv.Itoa(len(val))))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(_lineDelimiter))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(val))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(_lineDelimiter))
	return err
}

func decodeBulk(reader *bufio.Reader) (Bulk, error) {
	head, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	head = strings.TrimPrefix(head, string(_typeArrayPrefix))
	head = strings.TrimSuffix(head, _lineDelimiter)
	length, err := strconv.Atoi(head)
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return Bulk(""), nil
	}
	result := make(Bulk, length)
	_, err = io.ReadFull(reader, result)
	if err != nil {
		return Bulk(""), nil
	}
	reader.Discard(2)
	return result, nil
}

func encodeArray(w io.Writer, val Array) error {
	_, err := w.Write([]byte{_typeArrayPrefix})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(strconv.Itoa(len(val))))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(_lineDelimiter))
	if err != nil {
		return err
	}
	for _, item := range val {
		switch elem := item.(type) {
		case SimpleString:
			err = encodeSimpleString(w, elem)
		case Error:
			err = encodeError(w, elem)
		case Integer:
			err = encodeInteger(w, elem)
		case Bulk:
			err = encodeBulk(w, elem)
		case Array:
			err = encodeArray(w, elem)
		default:
			return errIllegalProto
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeArray(reader *bufio.Reader) (Array, error) {
	head, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	head = strings.TrimPrefix(head, string(_typeArrayPrefix))
	head = strings.TrimSuffix(head, _lineDelimiter)
	length, err := strconv.Atoi(head)
	if err != nil {
		return nil, err
	}
	body := make([]interface{}, 0, length)
	for idx := 0; idx < length; idx++ {
		_type, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		var val interface{}
		switch _type {
		case _typeSimplePrefix:
			val, err = decodeSimpleString(reader)
		case _typeErrorPrefix:
			val, err = decodeError(reader)
		case _typeIntegerPrefix:
			val, err = decodeInteger(reader)
		case _typeBulkPrefix:
			val, err = decodeBulk(reader)
		case _typeArrayPrefix:
			val, err = decodeArray(reader)
		default:
			return nil, errIllegalProto
		}
		if err != nil {
			return nil, err
		}
		body = append(body, val)
	}
	return body, nil
}
