package models

import "errors"

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Incident struct {
	Id     int
	Name   string
	Author int
}
type User struct {
	Id   int
	Name string
	Pass string
}
type UserRules struct {
	Id   int
	Name string
	User int
}
