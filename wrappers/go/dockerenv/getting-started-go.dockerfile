FROM ubuntu:16.04

RUN useradd -ms /bin/bash indy

# Install environment
RUN apt-get update -y && apt-get install -y \
	wget \
	python3.5 \
	python3-pip \
	python-setuptools \
	ipython \
	ipython-notebook \
	apt-transport-https \
	ca-certificates \
	software-properties-common

WORKDIR /home/indy

RUN pip3 install -U \
	pip \
	setuptools \
	jupyter \
	python3-indy==1.3.1-dev-476

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 68DB5E88 \
    && add-apt-repository "deb https://repo.sovrin.org/sdk/deb xenial master" \
    && apt-get update \
    && apt-get install -y \
    libindy

# Install Go
ENV GO_VERSION=1.9.1 \
    GOROOT=/goroot \
    GOPATH=/gopath 

ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN apt-get update -y && \
    apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git mercurial bzr && \
    mkdir /goroot && curl https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1 && \
    mkdir /gopath \
    mkdir /gopath/pkg

USER indy

EXPOSE 8888
