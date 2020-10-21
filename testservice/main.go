package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"time"
)

type Service struct {
	Name        string
	TTL         time.Duration
	ConsulAgent *consul.Agent
}

func NewService() (*Service, error) {
	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}

	s := &Service{
		Name:        "ht",
		TTL:         5 * time.Second,
		ConsulAgent: c.Agent(),
	}

	services, _, err := c.Health().Service("ht", "", true, nil)
	if err != nil {
		panic(err)
	}
	for _, a := range services {
		fmt.Println(a.Node.Address)
	}

	serviceDef := &consul.AgentServiceRegistration{
		Name: s.Name,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}

	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return nil, err
	}

	go s.UpdateTTL(s.Check)

	return s, nil
}

func (s *Service) check() bool {
	return true
}

func (s *Service) UpdateTTL(check func() (bool, error)) {
	ticker := time.NewTicker(s.TTL / 2)
	for range ticker.C {
		s.update(check)
	}
}

func (s *Service) update(check func() (bool, error)) {
	ok, err := check()
	if !ok {
		log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		if agentErr := s.ConsulAgent.FailTTL("service:"+s.Name, err.Error()); agentErr != nil {
			log.Print(agentErr)
		}
	} else {
		if agentErr := s.ConsulAgent.PassTTL("service:"+s.Name, ""); agentErr != nil {
			log.Print(agentErr)
		}
	}
}

func (s *Service) Check() (bool, error) {
	return true, nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	s, err := NewService()
	if err != nil {
		panic(err)
	}

	http.Handle("/", s)

	log.Fatal(http.ListenAndServe(":8001", nil))
}
