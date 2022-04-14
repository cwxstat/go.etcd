package etcdutils

import (
	"context"
	"fmt"
	"github.com/etcd-io/etcd/clientv3"
	"strconv"
	"testing"
	"time"
)

/*
To run tests, you must install KinD

Step 1:

curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.12.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind


Step 2:

kind create cluster


Step 3:

kubectl create configmap etcd-config --from-file=./certs/


Step 4:

kubectl apply -f pod.yaml


Step 5:

kubectl port-forward pods/etcd 2379:2379 -n default


Step 6:

go run main.go  config 


Step 7:
# From project root

etcdctl --endpoints=127.0.0.1:2379 --cacert=./certs/ca.pem \
 --cert=./certs/etcd-certs.pem --key=./certs/etcd-certs-key.pem user add root --password=A08auslkdjMMf 

 etcdctl --endpoints=127.0.0.1:2379 --cacert=./certs/ca.pem \
 --cert=./certs/etcd-certs.pem --key=./certs/etcd-certs-key.pem auth enable 

*/


func TestETC_Put(t *testing.T) {
	e, cancel := NewETC("test")
	defer cancel()

	e.DeleteWithPrefix("/testing")

	now := time.Now()
	e.Put("/testing/TestETC_Put", now.String())

	result, _ := e.Get("/testing/TestETC_Put")
	if string(result.Kvs[0].Value) != now.String() {
		t.Fatalf("%s\n", result.Kvs[0].Value)
	}

}

func TestETC_Delete(t *testing.T) {
	e, cancel := NewETC("test")
	defer cancel()

	e.DeleteWithPrefix("/testing")

	now := time.Now()

	e.Put("/testing/TestETC_Put ... more", now.String())
	e.PutWithLease("/testing/a", now.String(), 3)

	e.DeleteWithPrefix("/testing")
	result, _ := e.GetWithPrefix("/testing/")

	if len(result.Kvs) != 0 {
		t.Fatalf("Number of keys should be 0. You got: %d\n", len(result.Kvs))
	}

}

func TestETC_GetWithPrefix(t *testing.T) {
	e, cancel := NewETC("test")
	defer cancel()

	e.DeleteWithPrefix("/testing")

	now := time.Now()

	e.Put("/testing/TestETC_Put ... more", now.String())
	e.PutWithLease("/testing/a", now.String(), 3)

	result, _ := e.GetWithPrefix("/testing/")

	if len(result.Kvs) != 2 {
		t.Fatalf("Number of keys: %d\n", len(result.Kvs))
	}

	for i, v := range result.Kvs {
		sv := fmt.Sprintf("%s", v.Value)
		t.Logf("%s\n", sv)
		t.Logf("result.Kvs[%d]: %s, ver: %d,  lease: %d\n", i, v.Value, v.Version, v.Lease)
	}

}

/*
For this you neeed:
    "github.com/etcd-io/etcd/clientv3"
*/
func TestETC_Txn(t *testing.T) {
	e, cancel := NewETC("test")
	defer cancel()

	tx := e.Txn()

	txresp, err := tx.If(
		clientv3.Compare(clientv3.Value("foo"), "=", "bar"),
	).Then(
		clientv3.OpPut("foo", "sanfoo"), clientv3.OpPut("newfoo", "newbar"),
	).Else(
		clientv3.OpPut("foo", "bar"), clientv3.OpDelete("newfoo"),
	).Commit()
	fmt.Println(txresp, err)

	result, _ := e.Get("foo")
	for i, v := range result.Kvs {
		t.Logf("result.Kvs[%d]: %s, ver: %d,  lease: %d\n", i, v.Value, v.Version, v.Lease)
	}

}

func TestETC_Cli(t *testing.T) {
	e, cancel := NewETC("test")
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

func DoTxn(e ETC) {
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

func TestETC_Page(t *testing.T) {
	e, cancel := NewETC("test")
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

func TestETC_Options(t *testing.T) {
	e, cancel := NewETC("test")

	defer cancel()

	//e.DeleteWithPrefix("key")

	for i := 0; i < 20; i++ {
		k := "key"
		e.Put(k, strconv.Itoa(i))
	}

	opts := []clientv3.OpOption{
		//clientv3.AlarmMember{},
		clientv3.WithRev(605),
	}

	gr, err := e.Get("key", opts...)
	if err != nil {
		return
	}
	fmt.Println("--- First page ---")
	for _, item := range gr.Kvs {
		fmt.Println(string(item.Key), string(item.Value))
	}

}
