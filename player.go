package greatspacerace

import "net"

type Player struct {
	Name string
	Conn net.Conn
	Ship Ship
}
