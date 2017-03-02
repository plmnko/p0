// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

const sentBufferMaxSize = 100

type multiEchoServer struct {
	count        int                 // connected client
	ln           net.Listener        // listener for incoming connection
	connections  []*cliConnection    // keep track of connections
	mainExit     chan bool           // instruct serverChore to exit at Close()
	recBuf       chan string         // store msg from receive
	newCli       chan *cliConnection // new connection
	cliExit      chan int            // client's index for exit cleanup
	countRequest chan int            // for Count() request
}

type cliConnection struct {
	conn     net.Conn    // connection between server and client
	sentBuf  chan string // msg to client
	exitSend chan bool   // instruct send routine to exit
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	return &multiEchoServer{0, nil, nil, make(chan bool), make(chan string),
		make(chan *cliConnection), make(chan int), make(chan int)}

}

// should not block
func (mes *multiEchoServer) Start(port int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	mes.ln = ln
	fmt.Println("Listening at port: ", port)

	go mes.serverChores()
	go mes.handleConnection()

	return nil
}

func (mes *multiEchoServer) handleConnection() {
	for {
		conn, err := mes.ln.Accept()
		if err != nil {
			// throw errClosing at ln.Close()
			mes.ln.Close()
			return
		}
		// mes.countChange <- 1
		newCli := &cliConnection{conn, make(chan string, sentBufferMaxSize),
			make(chan bool)}
		mes.newCli <- newCli
		go mes.rec_data(newCli)
		go mes.send_data(newCli)
	}

}

func (mes *multiEchoServer) serverChores() {
	for {
		select {
		case <-mes.mainExit:
			return
		case <-mes.countRequest:
			mes.countRequest <- mes.count
		case newConn := <-mes.newCli:
			mes.count++
			mes.connections = append(mes.connections, newConn)
		case connIdx := <-mes.cliExit:
			mes.count--
			mes.connections[connIdx].exitSend <- true
			l := len(mes.connections)
			if connIdx == l-1 {
				mes.connections = mes.connections[:connIdx]
			} else {
				mes.connections = append(mes.connections[:connIdx],
					mes.connections[connIdx+1:]...)
			}
		case msg := <-mes.recBuf:
			for _, c := range mes.connections {
				// might block if c's buffer is full
				select {
				case c.sentBuf <- msg: // Put msg in the channel unless it is full
				default:
					// fmt.Printf("Channel for connection %p is full. len = %d\n", c, len(c.sentBuf))
				}
			}
		}
	}
}

func (mes *multiEchoServer) rec_data(cli *cliConnection) {
	// important! msg get lost if init reader for every msg
	reader := bufio.NewReader(cli.conn)
	for {

		line, err := reader.ReadString('\n')
		if err != nil {
			// disconnection
			for ix, c := range mes.connections {
				if cli == c {
					mes.cliExit <- ix
					break
				}
			}
			return
		} else {
			mes.recBuf <- line
		}

	}

}

func (mes *multiEchoServer) send_data(cli *cliConnection) {
	for {
		select {
		case <-cli.exitSend:
			return
		case msg := <-cli.sentBuf:
			_, err := cli.conn.Write([]byte(msg))
			if err != nil {
				return
			}
		}
	}
}

func (mes *multiEchoServer) Close() {

	if mes.ln != nil {
		mes.ln.Close()
	}

	for _, c := range mes.connections {
		c.exitSend <- true
		c.conn.Close() // cause receive to exit
	}
	mes.mainExit <- true

}

func (mes *multiEchoServer) Count() int {
	mes.countRequest <- 0
	return <-mes.countRequest
}
