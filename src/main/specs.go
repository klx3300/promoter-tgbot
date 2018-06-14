package main

import "time"

type Service struct {
	Name string
	Addr string
}

type Status struct {
	Tm         time.Time
	Stat       int
	Informator string
}

type Notification struct {
	Tm      time.Time
	Heading string
	Content string
}

type FetchRepsonse struct {
	Serv []Service
	Stat []Status
	Noti []Notification
}
