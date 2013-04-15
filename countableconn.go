package downloader

import (
	//"fmt"
	"net"
)

type CountableConnection struct {
	net.Conn
	WroteCounter   *IntCounter
	ReadCounter    *IntCounter
	ReadProgressCh chan int
}

func NewCountableConnection(conn net.Conn, progress chan int) CountableConnection {
	return CountableConnection{
		Conn:           conn,
		WroteCounter:   NewIntCounter(),
		ReadCounter:    NewIntCounter(),
		ReadProgressCh: progress,
	}
}

func (c CountableConnection) Write(p []byte) (n int, err error) {
	n, err = c.Conn.Write(p)
	c.WroteCounter.Incr(n)
	//fmt.Printf("CountableConnection.Write count = %d\n", c.WroteCounter.Value())
	return
}

func (c CountableConnection) Read(p []byte) (n int, err error) {
	n, err = c.Conn.Read(p)
	c.ReadCounter.Incr(n)
	//fmt.Printf("CountableConnection.Read count = %d\n", c.ReadCounter.Value())
	c.ReadProgressCh <- n
	return
}

/*
func (c CountableConnection) Close() (err error) {
	err = c.Conn.Close()
	fmt.Printf("CountableConnection.Close")
	return
}
*/
