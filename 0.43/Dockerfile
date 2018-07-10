FROM alpine:latest
MAINTAINER Juraj Bubniak <juraj.bubniak@gmail.com>

ENV HUGO_VERSION=0.43

RUN apk --no-cache add wget ca-certificates && \
  cd /tmp && \
  wget https://github.com/spf13/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz && \
  tar xzf hugo_${HUGO_VERSION}_Linux-64bit.tar.gz && \
  rm hugo_${HUGO_VERSION}_Linux-64bit.tar.gz && \
  mv hugo /usr/bin/hugo && \
  apk del wget ca-certificates

ENTRYPOINT ["/usr/bin/hugo"]