package resp

import "errors"

const _lineDelimiter = "\r\n"

const _typeSimplePrefix = '+'
const _typeErrorPrefix = '-'
const _typeIntegerPrefix = ':'
const _typeBulkPrefix = '$'
const _typeArrayPrefix = '*'

var errIllegalProto = errors.New("illegal RESP protocol.")
