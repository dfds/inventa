package misc

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	v1Networking "k8s.io/api/networking/v1beta1"
)

type InMemoryStore struct {
	Services  map[string]*v1.Service
	Ingresses map[string]*v1Networking.Ingress
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Services:  make(map[string]*v1.Service),
		Ingresses: make(map[string]*v1Networking.Ingress),
	}
}

// Service

func (s *InMemoryStore) PutService(key string, val *v1.Service) {
	s.Services[key] = val
	fmt.Println(string("\033[32m"), "Added: ", key)
	fmt.Println(string("\033[34m"), "Values: ", val, string("\033[0m"))
}

func (s *InMemoryStore) GetService(key string) *v1.Service {
	return s.Services[key]
}

func (s *InMemoryStore) GetServices() map[string]*v1.Service {
	return s.Services
}

func (s *InMemoryStore) RemoveService(key string) {
	delete(s.Services, key)
}

// Ingress

func (s *InMemoryStore) PutIngress(key string, val *v1Networking.Ingress) {
	s.Ingresses[key] = val
	fmt.Println(string("\033[32m"), "Added: ", key)
	fmt.Println(string("\033[34m"), "Values: ", val, string("\033[0m"))
}

func (s *InMemoryStore) GetIngress(key string) *v1Networking.Ingress {
	return s.Ingresses[key]
}

func (s *InMemoryStore) GetIngresses() map[string]*v1Networking.Ingress {
	return s.Ingresses
}

func (s *InMemoryStore) RemoveIngress(key string) {
	delete(s.Ingresses, key)
}
