ARG BASE_GO_IMAGE
FROM ${BASE_GO_IMAGE} AS base

RUN apk update && apk upgrade
RUN apk add --no-cache bash git openssh autoconf automake libtool gettext gettext-dev make g++ texinfo curl

# WORKDIR /root
# # fswatch is not available at alpine packages
# RUN wget https://github.com/emcrisostomo/fswatch/releases/download/1.14.0/fswatch-1.14.0.tar.gz
# RUN tar -xvzf fswatch-1.14.0.tar.gz
# WORKDIR /root/fswatch-1.14.0
# RUN ./configure
# RUN make && make install

ARG DEVELOPER_UID
RUN adduser -s /bin/sh -u ${DEVELOPER_UID} -D developer
USER developer

ARG APP_NAME
WORKDIR /go/src/${APP_NAME}
