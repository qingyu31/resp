package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

const _lineDelimiter = "\r\n"

const _typeSimplePrefix = '+'
const _typeErrorPrefix = '-'
const _typeIntegerPrefix = ':'
const _typeBulkPrefix = '$'
const _typeArrayPrefix = '*'

type Message interface {
	ReadFrom(reader *bufio.Reader) error
	WriteTo(writer *bufio.Writer) error
}

func WriteMessage(writer *bufio.Writer, message Message) error {
	return message.WriteTo(writer)
}

func ReadMessage(reader *bufio.Reader) (Message, error) {
	b, e := reader.ReadByte()
	if e != nil {
		return nil, e
	}
	var message Message
	switch rune(b) {
	case _typeSimplePrefix:
		message = new(SimpleStrings)
	case _typeIntegerPrefix:
		message = new(Integer)
	case _typeBulkPrefix:
		message = new(BulkStrings)
	case _typeArrayPrefix:
		message = new(Array)
	case _typeErrorPrefix:
		message = new(Error)
	default:
		return nil, errIllegalProto
	}
	e = message.ReadFrom(reader)
	return message, e
}

type SimpleStrings struct {
	value string
}

func NewSimpleStrings(val string) *SimpleStrings {
	s := new(SimpleStrings)
	s.value = val
	return s
}

func (s SimpleStrings) Value() string {
	return s.value
}

func (s *SimpleStrings) ReadFrom(reader *bufio.Reader) error {
	b, err := readLine(reader)
	if err != nil {
		return err
	}
	s.value = string(b)
	return nil
}

func (s SimpleStrings) WriteTo(writer *bufio.Writer) error {
	_, err := writer.WriteRune(_typeSimplePrefix)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(s.value)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(_lineDelimiter)
	return err
}

type Error struct {
	value error
}

func NewError(err error) *Error {
	e := new(Error)
	e.value = err
	return e
}

func (e Error) Value() error {
	return e.value
}

func (e *Error) ReadFrom(reader *bufio.Reader) error {
	b, err := readLine(reader)
	if err != nil {
		return err
	}
	e.value = errors.New(string(b))
	return nil
}

func (e Error) WriteTo(writer *bufio.Writer) error {
	_, err := writer.WriteRune(_typeErrorPrefix)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(e.value.Error())
	if err != nil {
		return err
	}
	_, err = writer.WriteString(_lineDelimiter)
	return err
}

type Integer struct {
	value int
}

func NewInteger(val int) *Integer {
	i := new(Integer)
	i.value = val
	return i
}

func (i Integer) Value() int {
	return i.value
}

func (i *Integer) ReadFrom(reader *bufio.Reader) error {
	b, err := readLine(reader)
	if err != nil {
		return err
	}
	i.value, err = strconv.Atoi(string(b))
	return err
}

func (i Integer) WriteTo(writer *bufio.Writer) error {
	_, err := writer.WriteRune(_typeIntegerPrefix)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(strconv.Itoa(i.value))
	if err != nil {
		return err
	}
	_, err = writer.WriteString(_lineDelimiter)
	return err
}

type BulkStrings struct {
	value []byte
}

func NewBulkStrings(val *string) *BulkStrings {
	bs := new(BulkStrings)
	if val != nil {
		bs.value = []byte(*val)
	}
	return bs
}

func NewBulkStringsWithBytes(val []byte) *BulkStrings {
	bs := new(BulkStrings)
	bs.value = val
	return bs
}

func (bs BulkStrings) Value() *string {
	if bs.value == nil {
		return nil
	}
	str := string(bs.value)
	return &str
}

func (bs *BulkStrings) ReadFrom(reader *bufio.Reader) error {
	b, err := readLine(reader)
	if err != nil {
		return err
	}
	length, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	if length < 0 {
		return nil
	}
	bs.value = make([]byte, length)
	_, err = io.ReadFull(reader, bs.value)
	if err != nil {
		return err
	}
	discardBytes := make([]byte, 2)
	_, err = io.ReadFull(reader, discardBytes)
	if err != nil {
		return err
	}
	return err
}

func (bs BulkStrings) WriteTo(writer *bufio.Writer) error {
	_, err := writer.WriteRune(_typeBulkPrefix)
	if err != nil {
		return err
	}
	if bs.value == nil {
		_, err = writer.WriteString("-1")
		if err != nil {
			return err
		}
		_, err = writer.WriteString(_lineDelimiter)
		return err
	}
	_, err = writer.WriteString(strconv.Itoa(len(bs.value)))
	if err != nil {
		return err
	}
	_, err = writer.WriteString(_lineDelimiter)
	if err != nil {
		return err
	}
	_, err = writer.Write(bs.value)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(_lineDelimiter)
	return err
}

type Array struct {
	value []Message
}

func NewArray(val []Message) *Array {
	a := new(Array)
	a.value = val
	return a
}

func (a Array) Value() []Message {
	return a.value
}

func (a *Array) ReadFrom(reader *bufio.Reader) error {
	b, err := readLine(reader)
	if err != nil {
		return err
	}
	length, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	if length < 0 {
		return nil
	}
	a.value = make([]Message, 0, length)
	for idx := 0; idx < length; idx++ {
		m, e := ReadMessage(reader)
		if e != nil {
			return e
		}
		a.value = append(a.value, m)
	}
	return nil
}

func (a Array) WriteTo(writer *bufio.Writer) error {
	_, err := writer.WriteRune(_typeArrayPrefix)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(strconv.Itoa(len(a.value)))
	if err != nil {
		return err
	}
	_, err = writer.WriteString(_lineDelimiter)
	if err != nil {
		return err
	}
	for _, msg := range a.value {
		err = msg.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	var buf bytes.Buffer
	var lastByte, currentByte byte
	var err error
	for {
		currentByte, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		buf.WriteByte(currentByte)
		if lastByte == '\r' && currentByte == '\n' {
			break
		}
		lastByte = currentByte
	}
	b := buf.Bytes()
	return b[0 : len(b)-2], nil
}
