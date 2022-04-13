package etcdutils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/cwxstat/go.etcd/pkg/settings"
	"io/ioutil"
	"log"
	"time"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
)

type ETC struct {
	CertsDir string
	ctx      context.Context
	cancel   context.CancelFunc
	Cli      *clientv3.Client
	kv       clientv3.KV
	err      error
	username string
	password string
}

func Server() settings.T {
	t := settings.T{}
	t.URL = "127.0.0.1:2379"
	t.TestURL = "127.0.0.1:2379"
	t.Certs.Directory = "/workspaces/go.etcd/certs"
	t.Certs.Ca = "ca.pem"
	t.Certs.Client = "etcd-certs.pem"
	t.Certs.ClientKey = "etcd-certs-key.pem"
	t.Username = "root"
	t.Password = "A08auslkdjMMf"
	return t
}

func NewETC(options ...string) (ETC, func()) {
	e := ETC{}
	var config settings.T
	var err error
	if options != nil && options[0] == "server" {

		config = Server()

	} else {

		config, err = settings.ReadConfig()
		if err != nil {
			log.Printf("You need a config. CREATING!")
			settings.CreateDefault()
			config, err = settings.ReadConfig()
			if err != nil {
				log.Fatalf("NewETC: Can't read or create config\n")
			}
		}
	}

	e.CertsDir = config.Certs.Directory
	e.username = config.Username
	e.password = config.Password
	url := config.URL
	if options != nil {
		if options[0] == "test" {
			url = config.TestURL
		}
	}

	e.ctx, e.cancel, e.Cli, e.kv, e.err = e.setup(config.Certs.Client,
		config.Certs.ClientKey, config.Certs.Ca, url)

	return e, e.cancel
}

func (e ETC) Cancel() {
	e.cancel()
	e.Cli.Close()
}

func (e ETC) setup(client, clientKey, ca, url string) (context.Context, context.CancelFunc, *clientv3.Client, clientv3.KV, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	cert, err := tls.LoadX509KeyPair(e.CertsDir+"/"+client, e.CertsDir+"/"+clientKey)
	caCert, err := ioutil.ReadFile(e.CertsDir + "/" + ca)
	caCertPool := x509.NewCertPool()

	if err != nil {
		return nil, nil, nil, nil, err
	}

	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, nil, nil, nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	//https://pkg.go.dev/go.etcd.io/etcd/clientv3
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: dialTimeout,
		Username:    e.username,
		Password:    e.password,
		TLS:         tlsConfig,
	})

	if err != nil {
		log.Fatalf("client3v3.New: %v\n", err)
	}
	kv := clientv3.NewKV(cli)
	return ctx, cancel, cli, kv, err
}

func (e ETC) Put(key string, value string) (*clientv3.PutResponse, error) {
	pr, err := e.kv.Put(e.ctx, key, value)
	return pr, err
}

func (e ETC) PutWithLease(key string, value string, ttl int64) (*clientv3.PutResponse, error) {
	lease, err := e.Cli.Grant(e.ctx, ttl)
	pr, err := e.kv.Put(e.ctx, key, value, clientv3.WithLease(lease.ID))
	return pr, err
}

func (e ETC) Get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	if opts != nil {
		gr, err := e.kv.Get(e.ctx, key, opts...)
		return gr, err
	}
	gr, err := e.kv.Get(e.ctx, key)
	return gr, err
}

func (e ETC) GetWithPrefix(key string) (*clientv3.GetResponse, error) {
	gr, err := e.kv.Get(e.ctx, key, clientv3.WithPrefix())
	return gr, err
}

func (e ETC) DeleteWithPrefix(key string) (*clientv3.DeleteResponse, error) {
	dr, err := e.kv.Delete(e.ctx, key, clientv3.WithPrefix())
	return dr, err
}

func (e ETC) Delete(key string) (*clientv3.DeleteResponse, error) {
	dr, err := e.kv.Delete(e.ctx, key)
	return dr, err

}

func (e ETC) Txn() clientv3.Txn {

	tx := e.kv.Txn(e.ctx)
	return tx
}
