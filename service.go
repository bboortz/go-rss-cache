package main

import (
	"time"
)

//type omit *struct{}

type Service struct {
	Id        int       `json:"id"`
	Name      string    `json:"Name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

type Services []Service
var services Services
var currentId = 0


// Give us some seed data¬                                                                                                                                                                                  
func init() {
	addService(Service{Name: "go-rnd"})
	addService(Service{Name: "go-keygen"})
}


func addService(s Service) Service {
	currentId += 1
	s.Id = currentId
	services = append(services ,s)
	logServiceRegistered(s)
	return s
}


