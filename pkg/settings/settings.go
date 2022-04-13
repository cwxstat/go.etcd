package settings

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var data = `
url: 127.0.0.1:2379
testurl: 127.0.0.1:2379
certs:
  directory: /workspaces/go.etcd/certs
  ca: ca.pem
  client: etcd-certs.pem
  clientKey: etcd-certs-key.pem
username: root
password: A08auslkdjMMf
tls: true
`

type T struct {
	URL     string
	TestURL string `yaml:"testurl"`
	Certs   struct {
		Directory string `yaml:"directory"`
		Ca        string `yaml:"ca"`
		Client    string `yaml:"client"`
		ClientKey string `yaml:"clientKey"`
	}
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	TLS      bool
}

func CreateDefault() {

	home, err := os.UserHomeDir()
	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile(home+"/"+".go.etcd.yaml", d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadConfig() (T, error) {
	t := T{}
	home, err := os.UserHomeDir()

	b, err := ioutil.ReadFile(home + "/" + ".go.etcd.yaml")
	if err != nil {
		return t, err
	}

	err = yaml.Unmarshal(b, &t)
	if err != nil {
		return t, err
	}
	return t, err

}

func TestRead() {
	t := T{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	d, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
