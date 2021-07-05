package common

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

var log = logrus.WithField("prefix", "registry")

// Service is a struct that con be registered into a ServiceRegistry
type Service interface {
	Start() error
	Stop() error
	Status() error
}

type ServiceRegistry struct {
	services 		map[reflect.Type]Service	// map of types to services
	serviceTypes 	[]reflect.Type				// keep an ordered slice of registered service types.
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[reflect.Type]Service),
	}
}

// StartAll initialized each service in order of registration.
func (s *ServiceRegistry) StartAll() {
	log.Debugf("Starting %d service: %v", len(s.serviceTypes), s.serviceTypes)
	for _, kind := range s.serviceTypes {
		log.Debugf("Starting service type %v", kind)
		go s.services[kind].Start()
	}
}

func (s *ServiceRegistry) StopAll() {
	for i:= len(s.serviceTypes) - 1; i >= 0; i-- {
		kind := s.serviceTypes[i]
		service := s.services[kind]
		if err := service.Stop(); err != nil {
			log.WithError(err).Errorf("Cound not stop the following service: %v", kind)
		}
	}
}

func (s *ServiceRegistry) Statuses() map[reflect.Type]error {
	m := make(map[reflect.Type]error, len(s.serviceTypes))
	for _, kind := range s.serviceTypes {
		m[kind] = s.services[kind].Status()
	}
	return m
}

// RegisterService appends a service constructor function to the service registry.
func (s *ServiceRegistry) RegisterService(service Service) error {
	kind := reflect.TypeOf(service)
	if _, exists := s.services[kind]; exists {
		return fmt.Errorf("servie already exists: %v", kind)
	}
	s.services[kind] = service
	s.serviceTypes = append(s.serviceTypes, kind)
	return nil
}

func (s *ServiceRegistry) FetchService(service interface{}) error {
	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		return fmt.Errorf("input must be of pointer type, received value type instead: %T", service)
	}
	element := reflect.ValueOf(service).Elem()
	if running, ok := s.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return fmt.Errorf("unknow service: %T", service)
}