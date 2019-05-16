FROM ubuntu:bionic

RUN apt-get update && apt-get install -y \
    ruby \
    ruby-dev \
    make \
    curl \
    wget \
    build-essential \
    git \
    rpm \
    zip \
    python \
    python-boto \
    python-jinja2

RUN gem install fpm --no-ri --no-rdoc

ENV GO_VERSION 1.10.3
ENV GO_ARCH amd64
ENV GOLANG go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
RUN wget -P /tmp https://dl.google.com/go/${GOLANG}; \
   tar -C /usr/local/ -xf /tmp/${GOLANG} ; \
   rm /tmp/${GOLANG}

RUN wget -O /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/ossia
ADD . /go/src/ossia
WORKDIR /go/src/ossia

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

VOLUME /go/src/ossia

ENTRYPOINT [ "/go/src/ossia/package.py" ]
