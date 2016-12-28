package main

import (
	"fmt"
	"time"
	"encoding/json"
)

type omit *struct{}

type Service struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

type Services []Service
var services Services
var currentId = 0


func addService(s Service) Service {
	currentId += 1
	s.Id = currentId
	services = append(services ,s)
	fmt.Println(s)
	return s
}


func (s *Service) getJson(s1 Service) string {
	b, err := json.Marshal(s1)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}

func (s *Service) getStruct(data []byte) Service {
	var s1 Service
	if err := json.Unmarshal(data, &s1); err != nil {
		return s1
	}
	return s1
}

func (s *Service) UnmarshalJSON(data []byte) error {
	var v [2]int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.Id= v[0]
	s.Name = string( v[1] )
	return nil
}
