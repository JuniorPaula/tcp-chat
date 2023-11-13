package main

import (
	"net"
	"log"
	"time"
	"fmt"
)

const (
	Port = "6969"
	SafeMode = true
	MessageRate = 1.0
	BanLimit = 10*60.0
	StrikeLimit = 10
)

func sensitive(message string) string {
	if SafeMode {
		return "[REDACTED]"
	} else {
		return message
	}
}

type MessageType int
const (
	ClientConnected MessageType = iota + 1
	ClientDisconnected
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

type Client struct {
	Conn net.Conn
	LastMesage time.Time
	StrikeCount int
}

func server(messages chan Message) {
	clients := map[string]*Client{}
	bannedMfs := map[string]time.Time{}

	for {
		msg := <- messages
		switch msg.Type {
		case ClientConnected:
			addr := msg.Conn.RemoteAddr().(*net.TCPAddr)
			bannedAt, banned := bannedMfs[addr.IP.String()]
			if banned {
				if time.Since(bannedAt).Seconds() >= BanLimit {
					delete(bannedMfs, addr.IP.String())
					banned = false
				}
			}
			if !banned {
				log.Printf("Client %s connected", sensitive(addr.String()))
				clients[msg.Conn.RemoteAddr().String()] = &Client{
					Conn: msg.Conn,
					LastMesage: time.Now(),
				}
			} else {
				msg.Conn.Write([]byte(fmt.Sprintf("You are banned MF: %f secs left\n", BanLimit - time.Since(bannedAt).Seconds())))
				msg.Conn.Close()
			}
		case ClientDisconnected:
			addr := msg.Conn.RemoteAddr().(*net.TCPAddr)
			log.Printf("Client %s connected", sensitive(addr.String()))
			delete(clients, addr.String())
		case NewMessage:
			authorAddr := msg.Conn.RemoteAddr().(*net.TCPAddr)
			author := clients[authorAddr.String()]
			now := time.Now()

			if author != nil {
				if now.Sub(author.LastMesage).Seconds() >= MessageRate {
					author.LastMesage = now
					author.StrikeCount = 0
					log.Printf("Client %s sent message %s", sensitive(authorAddr.String()), msg.Text)
					for _, client := range clients {
						if client.Conn.RemoteAddr().String() != authorAddr.String() {
							_, err := client.Conn.Write([]byte(msg.Text))
							if err != nil {
								// TODO: remove the connection from the list
								log.Printf("Could not send data to %s: %s\n", sensitive(client.Conn.RemoteAddr().String()), err)
							}
						}				
					}
				} else {
					author.StrikeCount++
					if author.StrikeCount >= StrikeLimit {
						bannedMfs[authorAddr.IP.String()] = now
						author.Conn.Write([]byte("You are banned!\n"))
						author.Conn.Close()
					}
				}
			} else {
				msg.Conn.Close()
			}
			
		}
	}
}

func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 64)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Cloud not read from %s: %s", sensitive(conn.RemoteAddr().String()), err)
			conn.Close()
			messages <- Message{
				Type: ClientDisconnected,
				Conn: conn,
			}
			return
		}

		text := string(buffer[0:n])
		messages <- Message{
			Type: NewMessage,
			Text: text,
			Conn: conn,
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("could not listen to epic port %s: %s\n", Port, err)
	}

	log.Printf("Listening to TCP connections on port %s ...\n", Port)

	messages := make(chan Message) 
	go server(messages)
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("could not accept a connection: %s\n", err)
		}
		log.Printf("Accepted connection from %s", sensitive(conn.RemoteAddr().String()))

		messages <- Message{
			Type: ClientConnected,
			Conn: conn,
		}
		
		go client(conn, messages)
	}

}
