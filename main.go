package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

type connection struct {
	Conn    net.Conn
	Index   int
	Message string
}

func (c *connection) Name() string {
	return fmt.Sprintf("Connection %d ", c.Index)
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}

// only needed below for sample processing
func getIP() net.IP {
	log.SetPrefix("IP ")
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			return ip // process IP address
		}
	}
	log.Fatal("Unable to find IP")
	return nil
}

func connect(done <-chan struct{}, ln net.Listener) <-chan connection {
	out := make(chan connection)
	go func(ln net.Listener) {
		defer close(out)
		keepRunning := true
		for i := 1; keepRunning; i++ {
			var err error
			conn := connection{Index: i}
			log.SetPrefix(fmt.Sprintf("Listener %d ", conn.Index))
			log.Println("on 54321")

			if conn.Conn, err = ln.Accept(); err != nil {
				log.Print(err)
			} else {
				log.SetPrefix(fmt.Sprintf("Listener %d ", conn.Index))
				log.Printf("Spawning Worker for %s", conn.Name())
			}
			select {
			case out <- conn:
			case <-done:
				keepRunning = false
				return
			}
		}
	}(ln)
	return out
}

func handle(done <-chan struct{}, in <-chan connection) <-chan connection {
	out := make(chan connection)
	go func() {
		for c := range in {
			log.SetPrefix(c.Name() + "Worker ")
			log.Println("Connected")

			for {
				var err error
				connReader := bufio.NewReader(c.Conn) //.ReadString('\n')

				log.SetPrefix(c.Name() + "Receive ")
				c.Message, err = connReader.ReadString('\n')
				if err != nil {
					if err.Error() == "EOF" {
						log.SetPrefix(c.Name() + "Worker ")
						log.Printf("%sClosed", c.Name())
						break
					}
					log.Fatal(err)
				}
				log.SetPrefix(c.Name() + "Receive ")
				log.Printf("%sMessage Received: %v", c.Name(), string(c.Message))
				select {
				case out <- c:
				case <-done:
					return
				}
			}
		}
		close(out)
	}()
	return out
}

func reply(in <-chan connection) {
	go func() {
		for c := range in {
			log.SetPrefix(c.Name() + "Reply ")
			c.Conn.Write([]byte("\n"))
		}
	}()
}

func main() {
	ip := getIP()

	log.SetPrefix("Init ")
	log.Println("Launching server...")
	log.Printf("Listening on %s:54321", ip.String())
	// listen on all interfaces
	ln, err := net.Listen("tcp", ":54321")

	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	conn := connect(done, ln)
	msg := handle(done, conn)
	reply(msg)

	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		ln.Close()
		done <- struct{}{}
		endWaiter.Done()
	}()
	endWaiter.Wait()
	//waitForCtrlC()

}
