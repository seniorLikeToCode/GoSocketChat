package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_MSG
	CMD_ROOMS
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

type room struct {
	name    string
	members map[net.Addr]*client
}

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	go c.readInput()
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. Usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. Usage: /join ROOM_NAME")
		return
	}

	roomName := args[1]
	r, ok := s.rooms[roomName]

	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required. Usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())
	s.quitCurrentRoom(c)
	c.msg("sad to see you go :(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(oldRoom.members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{id: CMD_NICK, client: c, args: args}
		case "/join":
			c.commands <- command{id: CMD_JOIN, client: c, args: args}
		case "/msg":
			c.commands <- command{id: CMD_MSG, client: c, args: args}
		case "/rooms":
			c.commands <- command{id: CMD_ROOMS, client: c}
		case "/quit":
			c.commands <- command{id: CMD_QUIT, client: c}
		default:
			c.msg(fmt.Sprintf("unknown command: %s", cmd))
		}
	}
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte(msg + "\n\n"))
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		}
	}
}

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error()) // Added error message
	}

	defer listener.Close()
	log.Printf("started server on 5000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
