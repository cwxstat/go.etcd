package main

import (
	"fmt"
	"github.com/cwxstat/go.etcd/pkg/etcdutils"
	"log"
	"time"
)

func Status() (string, error) {

	e, cancel := etcdutils.NewETC()
	defer cancel()

	result, err := e.GetWithPrefix("gopi.service")

	s := ""
	for i, v := range result.Kvs {
		s += fmt.Sprintf("result.Kvs[%d]: %s, ver: %d,  lease: %d\n", i, v.Value, v.Version, v.Lease)
	}
	return s, err
}

func Update() {

	e, cancel := etcdutils.NewETC("test")
	defer cancel()

	now := time.Now()

	e.PutWithLease("gopi.service/update", now.String(), 120)
	e.PutWithLease("gopi.service/addr", Ifconfig(), 120)

	result, _ := e.GetWithPrefix("gopi.service")

	for i, v := range result.Kvs {
		log.Printf("result.Kvs[%d]: %s, ver: %d,  lease: %d\n", i, v.Value, v.Version, v.Lease)
	}

}

func main() {

	for {
		Update()
		time.Sleep(70 * time.Second)
	}

}
