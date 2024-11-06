package socks4

import (
	errors "github.com/bdandy/go-errors"
)

const (
	socksVersion = 0x04
	socksConnect = 0x01

	accessGranted       = 0x5a
	accessRejected      = 0x5b
	accessIdentRequired = 0x5c
	accessIdentFailed   = 0x5d

	socksBind     = 0x02
	minRequestLen = 8
)

const (
	ErrorWrongNetWork = errors.String("network should be tcp or tcp4")
	ErrDialFailed     = errors.String("socks4 dial")
	ErrWrongAddr      = errors.String("wrong address: %s")

	ErrorHostUnknown   = errors.String("unable to find IP address of host %s")
	ErrBuffer          = errors.String("unable write into buffer")
	ErrIO              = errors.String("i\\o error")
	ErrIndentRequired  = errors.String("Ident required")
	ErrConnRejected    = errors.String("Connection rejected")
	ErrInvalidResponce = errors.String("Invalid response")
)

var Ident = "nobody@0.0.0.0"
