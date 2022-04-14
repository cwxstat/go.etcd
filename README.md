# go.etcd
etcd setup for Go, including backend

## Installation

```bash

# Step 1: 
docker run -it --rm -v ${PWD}:/work -w /work debian bash


# Step 2:
apt-get update && apt-get install -y curl &&
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_linux_amd64 -o /usr/local/bin/cfssl && \
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.5.0_linux_amd64 -o /usr/local/bin/cfssljson && \
chmod +x /usr/local/bin/cfssl && \
chmod +x /usr/local/bin/cfssljson

# Step 3:
mkdir -p ./certs
cfssl gencert -initca ./tls/ca-csr.json | cfssljson -bare ./certs/ca

# Step 4:
cfssl gencert \
  -ca=./certs/ca.pem \
  -ca-key=./certs/ca-key.pem \
  -config=./tls/ca-config.json \
  -hostname="etcd,etcd.default.svc.cluster.local,etcd.default.svc,localhost,etcd.pigbot.svc,127.0.0.1,etcd.pigbot.svc.cluster.local,34.111.92.27" \
  -profile=default \
  ./tls/ca-csr.json | cfssljson -bare ./certs/etcd-certs
```

# Download the etcd and build binary
```bash
git clone https://github.com/etcd-io/etcd.git
cd etcd
make
```

# Run etcd from root

```bash
# Here I'm assuming you're in codespaces
cd /workspaces/go.etcd
sudo chown -R codespace.codespace ./certs

./etcd/bin/etcd --client-cert-auth --trusted-ca-file=./certs/ca.pem \
--cert-file=./certs/etcd-certs.pem \
--key-file=./certs/etcd-certs-key.pem \
--peer-trusted-ca-file=./certs/ca.pem \
--advertise-client-urls=https://0.0.0.0:2379 --listen-client-urls=https://0.0.0.0:2379


```

# Test connection

Leave the above running and in a new terminal, run the following:


```bash
./etcd/bin/etcdctl --endpoints=127.0.0.1:2379 --cacert=./certs/ca.pem --cert=./certs/etcd-certs.pem --key=./certs/etcd-certs-key.pem put foo stuff

# and
./etcd/bin/etcdctl --endpoints=127.0.0.1:2379 --cacert=./certs/ca.pem --cert=./certs/etcd-certs.pem --key=./certs/etcd-certs-key.pem get foo 
```

# Set root password

```bash
# Note. This will prompt for the root password for etcd. Make up a password.  Here I'm using "A08auslkdjMMf
./etcd/bin/etcdctl --endpoints=127.0.0.1:2379 --cacert=./certs/ca.pem \
 --cert=./certs/etcd-certs.pem --key=./certs/etcd-certs-key.pem user add root --password=A08auslkdjMMf 


 ./etcd/bin/etcdctl --endpoints=127.0.0.1:2379 --cacert=./certs/ca.pem \
 --cert=./certs/etcd-certs.pem --key=./certs/etcd-certs-key.pem auth enable

```

# Now test connection again with root and password

```bash

./etcd/bin/etcdctl --endpoints=127.0.0.1:2379 \
--cacert=./certs/ca.pem --cert=./certs/etcd-certs.pem \
--key=./certs/etcd-certs-key.pem \
--user=root --password=A08auslkdjMMf get foo

```

# Steps to create configmap

```bash
kubectl create configmap etcd-config --from-file=./certs/

```


# Release Version

```bash
git tag -fa v0.0.5 -m "Update v0.0.5 tag"
git push origin v0.0.5 --force

```