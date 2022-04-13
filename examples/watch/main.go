package main

import (
	"context"
	"fmt"
	"github.com/cwxstat/go.etcd/pkg/etcdutils"
	"github.com/etcd-io/etcd/clientv3"
	"time"
)

func main() {

	e, cancel := etcdutils.NewETC("test")
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rch := e.Cli.Watch(ctx, "foo", clientv3.WithPrefix())

	go func(chn clientv3.WatchChan) {
		for wresp := range chn {
			for _, ev := range wresp.Events {
				fmt.Printf("WATCH!!")
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}(rch)

	for i := 0; i < 12; i++ {
		DoTxn(e)
		time.Sleep(300 * time.Millisecond)
	}

}

func DoTxn(e etcdutils.ETC) {
	tx := e.Txn()

	_, err := tx.If(
		clientv3.Compare(clientv3.Value("foo"), "=", "bar"),
	).Then(
		clientv3.OpPut("foo", "sanfoo"), clientv3.OpPut("newfoo", "newbar"),
	).Else(
		clientv3.OpPut("foo", "bar"), clientv3.OpDelete("newfoo"),
	).Commit()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(txresp, err)
}
