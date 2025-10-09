package service

import (
	"bytes"
	"fmt"
	"net/url"
	"sync"
)

type ServiceMap interface {
	// Do not include slashes in the prefix.
	AddEndpoint(rawUrl, prefix, scope string, public bool)
	GetEndpoint(prefix string) (Endpoint, error)
	PrintAll()
}

type serviceMap struct {
	mtx      sync.RWMutex
	services map[string]Endpoint
}

type Endpoint struct {
	Url    *url.URL
	Prefix string
	Scope  string
	Public bool
}

var serviceMapInstance *serviceMap

func NewServiceMap() ServiceMap {
	if serviceMapInstance != nil {
		return serviceMapInstance
	}

	newServiceMap := &serviceMap{
		mtx:      sync.RWMutex{},
		services: make(map[string]Endpoint),
	}
	serviceMapInstance = newServiceMap

	return serviceMapInstance
}

func (s *serviceMap) AddEndpoint(rawUrl, prefix, scope string, public bool) {
	// TODO: more checks
	url, err := url.Parse(rawUrl)
	if err != nil {
		return
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, ok := s.services[prefix]
	if ok {
		return
	}

	s.services[prefix] = Endpoint{
		Url:    url,
		Prefix: prefix,
		Scope:  scope,
		Public: public,
	}
}

func (s *serviceMap) GetEndpoint(prefix string) (Endpoint, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	endpoint, ok := s.services[prefix]
	if !ok {
		return Endpoint{}, fmt.Errorf("error: does not exist")
	}

	return endpoint, nil
}

func (s *serviceMap) PrintAll() {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	fmt.Println("======SERVICE MAP======")
	for _, v := range s.services {
		fmt.Printf("%v\n", v.String())
	}
	fmt.Println("======END======")
}

func (e *Endpoint) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("Endpoint:\n"))
	buf.WriteString(fmt.Sprintf("Url: '%v'\n", e.Url))
	buf.WriteString(fmt.Sprintf("Prefix: '%v'\n", e.Prefix))
	buf.WriteString(fmt.Sprintf("Scope: '%v'\n", e.Scope))
	buf.WriteString(fmt.Sprintf("Public: '%v'\n", e.Public))

	return buf.String()
}
