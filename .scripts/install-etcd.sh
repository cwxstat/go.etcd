#!/bin/bash
function is_bin_in_path {
  builtin type -P "$1" &> /dev/null
}

is_bin_in_path ectdctl && echo "ectdctl is already installed" && exit 0

ETCD_VER=v3.5.3

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}

rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

/tmp/etcd-download-test/etcd --version
/tmp/etcd-download-test/etcdctl version
/tmp/etcd-download-test/etcdutl version

mv /tmp/etcd-download-test/etcdctl /home/codespace/.local/bin/ectdctl
mv /tmp/etcd-download-test/etcd  /home/codespace/.local/bin/etcdctl
mv /tmp/etcd-download-test/etcdutl  /home/codespace/.local/bin/etcdutl
