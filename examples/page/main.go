package main

import (
	"fmt"
	"github.com/cwxstat/go.etcd/pkg/etcdutils"
	"github.com/etcd-io/etcd/clientv3"
	"strconv"
)

func main() {

	e, cancel := etcdutils.NewETC("test")
	defer cancel()

	e.DeleteWithPrefix("key")

	for i := 0; i < 20; i++ {
		k := fmt.Sprintf("key_%02d", i)
		e.Put(k, strconv.Itoa(i))
	}

	var number int64 = 3
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(number),
	}

	gr, _ := e.Get("key", opts...)
	fmt.Println("--- First page ---")
	for _, item := range gr.Kvs {
		fmt.Println(string(item.Key), string(item.Value))
	}

	lastKey := string(gr.Kvs[len(gr.Kvs)-1].Key)

	fmt.Println("--- Second page ---")
	opts[2] = clientv3.WithLimit(number + 1)
	opts = append(opts, clientv3.WithFromKey())
	gr, _ = e.Get(lastKey, opts...)

	// Skipping the first item, which the last item from from the previous Get
	for _, item := range gr.Kvs[1:] {
		fmt.Println(string(item.Key), string(item.Value))
	}

}
