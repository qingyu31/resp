package resp

import "errors"

var errIllegalProto = errors.New("illegal RESP protocol.")
var errIllegalRequest = errors.New("illegal RESP request.")
var errIllegalResponse = errors.New("illegal RESP response.")
