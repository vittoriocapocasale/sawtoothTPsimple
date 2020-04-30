
FROM ubuntu:bionic

RUN apt-get update \
 && apt-get install gnupg -y

LABEL "install-type"="mounted"

RUN echo "deb [arch=amd64] http://repo.sawtooth.me/ubuntu/ci bionic universe" >> /etc/apt/sources.list \
 && (apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD \
 || apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 8AA7AF1F1091A5FD) \
 && apt-get update \
 && apt-get install -y -q \
    golang-1.10-go \
    git \
    libssl-dev \
    libzmq3-dev \
    openssl \
    dpkg-dev \
    protobuf-compiler \
    python3 \
    python3-grpcio \
    python3-grpcio-tools \
    python3-pkg-resources \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

ENV GOPATH=/go

ENV PATH=$PATH:/go/bin:/usr/lib/go-1.10/bin

RUN mkdir /go

RUN go get -u  \
    github.com/btcsuite/btcd/btcec \
    github.com/golang/protobuf/proto \
    github.com/golang/protobuf/protoc-gen-go \
    github.com/golang/mock/gomock \
    github.com/golang/mock/mockgen \
    github.com/ugorji/go/codec \
    github.com/jessevdk/go-flags \
    github.com/pebbe/zmq4 \
    github.com/brianolson/cbor_go \
    github.com/pelletier/go-toml \
    golang.org/x/crypto/ripemd160 \
    golang.org/x/crypto/ssh \
    github.com/hyperledger/sawtooth-sdk-go \
    github.com/satori/go.uuid 
    
RUN mkdir -p /go/src/github.com/vittoriocapocasale/easy_tp

WORKDIR /go/src/github.com/vittoriocapocasale/easy_tp/

COPY ./ /go/src/github.com/vittoriocapocasale/easy_tp/

RUN mkdir -p /go/src/github.com/hyperledger/sawtooth-sdk-go 

WORKDIR /go/src/github.com/hyperledger/sawtooth-sdk-go 

RUN go generate

WORKDIR /go/src/github.com/vittoriocapocasale/easy_tp

RUN go install

EXPOSE 4004/tcp

CMD ["easy_tp", "-C", "tcp://validator:4004"]
