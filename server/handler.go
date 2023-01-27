package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// функция handleUser обрабатывает всех входящих пользователей, которые заходят в чат. Проверяет, существует ли уже имя пользователя
// в чате. Транслируйте всем существующим пользователям о входе нового пользователя и отправьте сообщение от имени пользователя.

func (c *Client) handleUsers(conn net.Conn) {
	var userName string
	reader := bufio.NewReader(conn)
	userName, err := reader.ReadString('\n')
	if err != nil {
		l.Printf("error with username: %v\n", err)
	}
	userName = strings.Trim(userName, " \r\n	")
	for !c.checkName(userName) {
		conn.Write([]byte("User name is already taken or written incorrectly. \n[choose another one]:"))
		userName, _ = reader.ReadString('\n')
		userName = strings.Trim(userName, " \r\n	")
	}

	defer conn.Close()

	c.Lock()
	c.Users[conn] = userName
	c.Unlock()

	conn.Write([]byte(c.History))

	c.Lock()
	c.LeaveJoin <- c.notification("\n"+userName+" has joined the chat...\n", conn)
	c.Unlock()
	l.Println(userName + " has joined the chat..") // дублирование в журнал сервера
	defer conn.Close()

	for {
		var exit bool
		var input string
		var err error
		time := time.Now().Format("2006-01-02 15:04:05")
		conn.Write([]byte(fmt.Sprintf("[%s] [%s]:", time, userName)))
		for {
			input, err = reader.ReadString('\n')
			if err != nil {
				exit = true
				break
			} else if !validMsg(input) {
				conn.Write([]byte(brackets(time) + " " + brackets(userName) + ":"))
			} else {
				break
			}
		}
		if exit {
			c.Lock()
			delete(c.Users, conn) // удаляет уходящего пользователя из карты пользователей
			c.LeaveJoin <- c.notification("\n"+userName+" has left the chat...\n", conn)
			c.Unlock()
			l.Println(userName + " has left the chat...") // duplicting to server log
			break
		}
		c.Lock()
		c.Messages <- c.newMessage(brackets(userName)+":"+input, conn)
		c.Unlock()
		defer conn.Close()
	}
}
