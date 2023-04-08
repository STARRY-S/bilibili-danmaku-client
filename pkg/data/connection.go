package data

import "io"

type WsConnection interface {
	io.ReadWriteCloser

	IsClientConn() bool
	IsServerConn() bool
}
