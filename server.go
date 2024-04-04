package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CONN_PORT = ":4445"
	CONN_TYPE = "tcp"
)

var conn_cnt int = 0

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	log.Println("Listening on" + CONN_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			continue
		}
		conn_cnt++
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		len, err := conn.Read(buffer)
		if err != nil {
			log.Println("Client disconnected or error reading:", err.Error())
			return
		}
		receivedText := strings.TrimSpace(string(buffer[:len]))

		if receivedText == "/join" {
			conn.Write([]byte("Welcome to the chat!\n"))
		} else if receivedText == "/view" {
			conn.Write([]byte("The current connection cnt is: " + strconv.Itoa(conn_cnt) + "\n"))
		} else {
			conn.Write([]byte("Message received: " + receivedText + "\n"))
		}

		logMessage(receivedText)
	}
}

func logMessage(message string) {
	f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s - %s\n", timestamp, message)
	if _, err := f.WriteString(logEntry); err != nil {
		log.Println(err)
	}
}
