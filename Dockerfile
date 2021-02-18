# This is the toolbox dockerfile that provides all necessary tools to build portainer
FROM ubuntu:focal-20210119

# Set TERM as noninteractive to suppress debconf errors
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

# Set default go version
ARG GO_VERSION=go1.13.11.linux-amd64

# Install packages
RUN apt-get update --fix-missing && apt-get install -qq \
    dialog \
    apt-utils \
    curl \
    build-essential \
    nodejs \
    git \
    wget

# Install Yarn
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
    && echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
    && apt-get update && apt-get -y install yarn

# Install Golang
RUN cd /tmp \
    && wget -q https://dl.google.com/go/${GO_VERSION}.tar.gz \
    && tar -xf ${GO_VERSION}.tar.gz \
    && mv go /usr/local

# Configure Go
ENV PATH "$PATH:/usr/local/go/bin"