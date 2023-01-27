package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	outfile, _ = os.Create("log.log")    // создание файла для журналов
	l          = log.New(outfile, "", 0) // все журналы будут сохранены в файл
	maxUsers   = 10                      // максимально допустимое количество пользователей
)

type Client struct {
	sync.Mutex
	Users     map[net.Conn]string
	LeaveJoin chan Message
	Messages  chan Message
	History   []byte
}

type Message struct {
	Text    string
	Address string
	Time    string
}

func NewServer() *Client {
	return &Client{
		Users:     make(map[net.Conn]string),
		LeaveJoin: make(chan Message),
		Messages:  make(chan Message),
	}
}

// Функция запуска запускает сервер с заданным портом..
func (c *Client) Run(port string) {
	p := fmt.Sprintf(":%s", port)
	listen, err := net.Listen("tcp", p) // инициализировать сервер
	if err != nil {
		l.Printf("server error: %v\n", err.Error())
		return

	}
	log.Printf("Server is listening on port %s\n", port) // для отображения в терминале
	l.Printf("Server is listening port %s\n", port)      // для сохранения в журнале mylog.log

	welcome := Logo() // получение логотипа для входа в чат
	defer listen.Close()

	go c.broadcaster() // функция объявления, принимающая входящие данные канала
	for {
		conn, err := listen.Accept() // прослушивание входящих соединений
		if err != nil {
			l.Printf("error with listen conn: %v\n", err)
			return
			// continue
		}
		c.Lock()
		if len(c.Users) >= maxUsers {
			l.Println("server is full")
			conn.Write([]byte("Server is already full. Please try it later.\n"))
			conn.Close()
		} else {
			conn.Write(welcome)
			conn.Write([]byte("\n[ENTER YOUR NAME]:"))
			go c.handleUsers(conn)
		}
		c.Unlock()
	}
}

// функция newMessage собирает данные и помещает их в strtuct, который будет отправлен в канал.
// Дополнительно struct заполняется адресом пользователя, текущим временем и сообщением пользователя.
// Также сохраняет само текущее сообщение в файл и распечатывает его со стороны сервера. Также сохраняются все логи
// в дополнительный файл.
func (c *Client) newMessage(msg string, conn net.Conn) Message {
	addr := conn.RemoteAddr().String()
	time := time.Now().Format("2006-01-02 15:04:05")
	c.History = append(c.History, []byte(brackets(time)+" "+msg)...)
	l.Println(brackets(time) + " " + msg) // дублирование в журнал сервера
	os.WriteFile("./history.txt", c.History, 0o666)
	return Message{
		Text:    msg,
		Address: addr,
		Time:    time,
	}
}

func (c *Client) notification(msg string, conn net.Conn) Message {
	addr := conn.RemoteAddr().String()
	time := time.Now().Format("2006-01-02 15:04:05")
	l.Println(msg) // дублирование в журнал сервера
	os.WriteFile("./history.txt", c.History, 0o666)
	return Message{
		Text:    msg,
		Address: addr,
		Time:    time,
	}
}

// Функция broadcaster - это goroutine, которая принимает сигналы каналов.
func (c *Client) broadcaster() {
	for {
		select {
		case msg := <-c.Messages:
			c.Lock()
			for conn, user := range c.Users {
				if msg.Address == conn.RemoteAddr().String() {
					continue
				}
				conn.Write([]byte(fmt.Sprintf("\n[%s] %s", msg.Time, msg.Text)))
				conn.Write([]byte(fmt.Sprintf("[%s] [%s]:", msg.Time, user))) // иметь постоянно отображаемую строку для ввода данных пользователем
			}
			c.Unlock()
		case msg := <-c.LeaveJoin:
			c.Lock()
			for conn, user := range c.Users {
				if msg.Address == conn.RemoteAddr().String() {
					continue
				}
				conn.Write([]byte(msg.Text))
				conn.Write([]byte(fmt.Sprintf("[%s] [%s]:", msg.Time, user)))
			}
			c.Unlock()
		}
	}
}
