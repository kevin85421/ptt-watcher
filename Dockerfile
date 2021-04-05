FROM ubuntu:18.04
MAINTAINER Apache Software Foundation <dev@submarine.apache.org>

RUN apt-get update &&\
    apt-get install -y wget vim git

ENV GOROOT="/usr/local/go"
ENV GOPATH=$HOME/gocode
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOPATH:$GOBIN:$GOROOT/bin

RUN wget https://dl.google.com/go/go1.16.2.linux-amd64.tar.gz &&\
    tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz

RUN cd /usr/src &&\
    git clone https://github.com/kevin85421/ptt-watcher.git &&\
    cd ptt-watcher &&\
    go mod vendor &&\
    go build -o ptt-watcher &&\
    chmod +x ptt-watcher

ADD config.json /usr/src/ptt-watcher/config.json
WORKDIR /usr/src/ptt-watcher
CMD ["./ptt-watcher"]~
