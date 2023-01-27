package server

import (
	"fmt"
	"log"
	"os"
)

// Функция Logo считывает файл txt с логотипом пингвина и сохраняет его как массив байтов.
func Logo() []byte {
	logo, err := os.ReadFile("./welcome.txt")
	if err != nil {
		log.Printf("logo error %s\n", err)
	}
	return logo
}

// Функция bracket принимает строку с содержимым и возвращает отформатированное содержимое со скобками.
func brackets(s string) string {
	return "[" + s + "]"
}

// Функция validMsg проверяет, является ли само содержимое пробелом или новой строкой, и возвращает булево значение в соответствии с этим.
func validMsg(msg string) bool {
	for i := range msg {
		if msg[i] != ' ' && msg[i] != '\n' {
			return true
		}
	}
	return false
}

// Функция checkName проверяет, был ли уже введен в чат пользователь с таким же именем.
// Также если имя пользователя правильное.
func (c *Client) checkName(s string) bool {
	if !validMsg(s) {
		return false
	}
	for _, user := range c.Users {
		if s == user {
			return false
		}
	}
	return true
}

// Функция PortCheck проверяет правильность ввода в поле порта. Оно должно содержать 4 цифры.
func PortCheck(s string) bool {
	if len(s) != 4 {
		fmt.Println("Incorrect port length")
		return false
	}
	for _, w := range s {
		if w < '0' || w > '9' {
			return false
		}
	}
	return true
}
