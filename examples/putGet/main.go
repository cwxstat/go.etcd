package main

import (
	"github.com/cwxstat/go.etcd/pkg/etcdutils"
	"log"
	"time"
)

func main() {

	e, cancel := etcdutils.NewETC()
	defer cancel()

	e.DeleteWithPrefix("/testing")

	now := time.Now()

	e.Put("/testing/TestETC_Put ... more", now.String())
	e.PutWithLease("/testing/a", now.String(), 3)

	result, _ := e.GetWithPrefix("/testing/")

	if len(result.Kvs) != 2 {
		log.Fatalf("Number of keys: %d\n", len(result.Kvs))
	}

	for i, v := range result.Kvs {
		log.Printf("result.Kvs[%d]: %s, ver: %d,  lease: %d\n", i, v.Value, v.Version, v.Lease)
	}

}
