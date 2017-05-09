FROM alpine:latest
MAINTAINER Juraj Bubniak <juraj.bubniak@gmail.com>

ENV HUGO_VERSION=0.16

RUN apk --no-cache add wget ca-certificates && \
  cd /tmp/ && \
  wget https://github.com/spf13/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-64bit.tgz && \
  tar xzf hugo_${HUGO_VERSION}_linux-64bit.tgz && \
  rm hugo_${HUGO_VERSION}_linux-64bit.tgz && \
  mv hugo /usr/bin/hugo && \
  apk del wget ca-certificates

ENTRYPOINT [ "hugo" ]