package resp

import "v9.git.n.xiaomi.com/miot_shop/go_vendors/github.com/go-errors/errors"

const LINE_DELIMETER = "\r\n"

var errIllegalProto = errors.New("illegal RESP protocol.")
