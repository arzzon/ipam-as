package main

import (
	"fmt"
	"github/arzzon/ipam-as/pkg/dummy"
	"time"
)

func main() {
	var ds *dummy.DummyStruct
	ds = dummy.NewDummyStruct("dm1", 1)
	fmt.Printf("Successfully created dummy obj: Name: %s, ID:%d.", ds.GetName(), ds.GetID())
	time.Sleep(3000 * time.Second)
}
