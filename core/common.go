package core

import (
	"io"
	"log"
	"net"
	"sync"
)

// 重连间隔时间
const redialIntervalTime = 5

func connCopy(source, target net.Conn, wg *sync.WaitGroup) {
	_, err := io.Copy(source, target)
	if err != nil {
		log.Println("Connection interrupted", err.Error())
	}
	wg.Done()
}

func forward(conn1, conn2 net.Conn) {
	log.Println(&conn1)
	log.Printf("Forward channel [%s/%s] <-> [%s/%s]\n",
		conn1.RemoteAddr(), conn1.LocalAddr(), conn2.RemoteAddr(), conn2.LocalAddr())

	var wg sync.WaitGroup
	// wait tow goroutines
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	//blocking when the wg is locked
	wg.Wait()
}
