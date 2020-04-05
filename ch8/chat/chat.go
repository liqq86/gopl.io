// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//!+broadcaster
type clientID string
type cmsg struct {
	cid   clientID
	msg   string
	mtype int
}

type client chan<- cmsg // an outgoing message channel
type clientEntry struct {
	cid clientID
	ch  client
}

var (
	entering = make(chan clientEntry)
	leaving  = make(chan clientEntry)
	messages = make(chan cmsg) // all incoming client messages
)

func broadcaster() {
	clients := make(map[clientID]client) // all connected clients
	for {
		select {
		case m := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			if m.mtype == 2 {
				log.Println("say " + m.msg + " " + string(m.cid))
			}
			for cid, cli := range clients {
				if m.cid != cid {
					cli <- m
				}
			}

		case ce := <-entering:
			log.Println("enter " + ce.cid)
			clients[ce.cid] = ce.ch

		case ce := <-leaving:
			log.Println("leave " + ce.cid)
			delete(clients, ce.cid)
			close(ce.ch)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan cmsg) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	cid := clientID(who)
	ch <- cmsg{cid, "You are " + who, 1}
	messages <- cmsg{cid, who + " has arrived", 0}
	entering <- clientEntry{clientID(who), ch}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- cmsg{cid, input.Text(), 2}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- clientEntry{clientID(who), ch}
	messages <- cmsg{cid, who + " has left", 0}
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan cmsg) {
	for msg := range ch {
		fmt.Fprintln(conn, msg.msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
