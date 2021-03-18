package message

import (
	"errors"
	"sync"
)

type Property interface {
	Get(key string) (string, error)
	Set(key string, value string)
	Remove(key string) error
	Clear()
	SetMap(propertyMap map[string]string)
	GetMap() map[string]string
}

type IProperty struct {
	proLock sync.RWMutex
	data map[string]string
}

func (p *IProperty) Get(key string) (string, error) {
	p.proLock.RLock()
	defer p.proLock.RUnlock()
	if value, ok := p.data[key]; ok {
		return value, nil
	} else {
		return "", errors.New("key does not exist")

	}
}
func (p *IProperty) Set(key string, value string) {
	p.proLock.Lock()
	defer p.proLock.Unlock()
	p.data[key] = value

}
func (p *IProperty) Remove(key string) error {
	p.proLock.Lock()
	defer p.proLock.Unlock()
	if _, ok := p.data[key]; ok {
		delete(p.data, key)
		return nil
	}
	return errors.New("key does not exist")
}
func (p *IProperty) GetMap() map[string]string{

	return p.data
}
func (p *IProperty) SetMap(propertyMap map[string]string) {
	p.proLock.Lock()
	defer p.proLock.Unlock()
	p.data = propertyMap
}


func (p *IProperty) Clear() {
	p.proLock.Lock()
	defer p.proLock.Unlock()

	for k, _ := range p.data {
		delete(p.data, k)
	}

}
