package dummy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type DummyStruct struct {
	Name string
	ID   int
}

func NewDummyStruct(name string, ID int) *DummyStruct {
	log.WithFields(log.Fields{
		"level": "INFO",
	}).Info("Creating DummyStruct")
	fmt.Println("Creating struct")
	return &DummyStruct{Name: name, ID: ID}
}

func (ds *DummyStruct) GetID() int {
	return ds.ID
}

func (ds *DummyStruct) GetName() string {
	return ds.Name
}
