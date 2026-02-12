FROM golang:latest

WORKDIR /app

# Go araçlarını yükle
RUN go install \
    golang.org/x/tools/gopls@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest


RUN apt-get update && apt-get install -y \
    git \
    curl \
    wget \
    build-essential \
    cmake \
    unzip \
    gettext \
 #   python3 \
  #  python3-pip \
  #  nodejs \
  #  npm \
    ripgrep \
    fd-find \
    && rm -rf /var/lib/apt/lists/*

#RUN curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux64.tar.gz && \
#   tar -xzf nvim-linux64.tar.gz && \
#    mv nvim-linux64 /usr/local/nvim && \
#    ln -s /usr/local/nvim/bin/nvim /usr/local/bin/nvim && \
#    rm nvim-linux64.tar.gz

RUN curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-arm64.appimage && \
    chmod u+x nvim-linux-arm64.appimage && \
    ./nvim-linux-arm64.appimage --appimage-extract && \
    mv squashfs-root /usr/local/nvim && \
    ln -s /usr/local/nvim/AppRun /usr/local/bin/nvim && \
    rm nvim-linux-arm64.appimage

# Environment variables
ENV GO111MODULE=on \
    GOPATH=/go \
    PATH=$PATH:/go/bin
