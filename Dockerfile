# Copyright 2018 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -----------------------------------------------------------------------------

FROM ubuntu:bionic

ENV VERSION=AUTO_STRICT

RUN apt-get update \
 && apt-get install gnupg -y

RUN echo "deb [arch=amd64] http://repo.sawtooth.me/ubuntu/ci bionic universe" >> /etc/apt/sources.list \
 && echo "deb http://archive.ubuntu.com/ubuntu bionic-backports universe" >> /etc/apt/sources.list \
 && echo 'deb http://ppa.launchpad.net/gophers/archive/ubuntu bionic main' >> /etc/apt/sources.list \
 && (apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD \
 || apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 8AA7AF1F1091A5FD) \
 && (apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 308C15A29AD198E9 \
 || apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 308C15A29AD198E9) \
 && apt-get update \
 && apt-get install -y -q \
    build-essential \
    golang-1.13-go \
    git \
    libssl-dev \
    libzmq3-dev \
    openssl \
    python3-grpcio-tools \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

#RUN mkdir -p /auctioin
ENV GOPATH=/go:/go/src/github.com/hyperledger/sawtooth-sdk-go:/go/src/github.com/moorada/sawtooth-auction-go
ENV PATH=$PATH:/project/bin:/go/bin:/usr/lib/go-1.13/bin

RUN mkdir /go

RUN go get -u \
    github.com/golang/protobuf/proto \
    github.com/golang/protobuf/protoc-gen-go \
    github.com/pebbe/zmq4 \
    github.com/satori/go.uuid \
    github.com/btcsuite/btcd/btcec \
    github.com/jessevdk/go-flags \
    github.com/golang/mock/gomock \
    github.com/golang/mock/mockgen \
    golang.org/x/crypto/ssh \
    github.com/hyperledger/sawtooth-sdk-go \
    github.com/moorada/sawtooth-auction-go/src/sawtooth_auction

WORKDIR /go/src/github.com/hyperledger/sawtooth-sdk-go
RUN go generate

EXPOSE 4004/tcp

WORKDIR /go/src/github.com/moorada/sawtooth-auction-go

RUN bash -c "./install.sh"
#CMD ["auction", "-vv"]
