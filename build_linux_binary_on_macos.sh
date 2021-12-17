#!/bin/bash
# build_linux_binary_on_macos.sh

# NOTE: you need install follows on Linux.
# wget https://copr.fedorainfracloud.org/coprs/ngompa/musl-libc/repo/epel-7/ngompa-musl-libc-epel-7.repo -O /etc/yum.repos.d/ngompa-musl-libc-epel-7.repo --no-check-certificate
# yum install -y musl-libc-static

brew install FiloSottile/musl-cross/musl-cross

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -o traffic-clones-api_linux_amd64

exit 0
