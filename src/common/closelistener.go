package common

import "net"

type CloseListener interface {
	OnErrorHeader(conn *net.TCPConn)
	OnConnectClose(conn *net.TCPConn)
}
