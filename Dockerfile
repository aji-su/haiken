FROM golang:1.17-buster

RUN apt-get update \
  && apt-get install --no-install-recommends -y \
    curl \
    file \
    git \
    libmecab-dev \
    make \
    mecab \
    mecab-ipadic-utf8 \
    patch \
    sudo \
    xz-utils \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /opt
RUN git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git

WORKDIR /opt/mecab-ipadic-neologd
RUN ./bin/install-mecab-ipadic-neologd -n -y \
  && rm -rf /opt/mecab-ipadic-neologd \
  && echo "dicdir = /usr/lib/x86_64-linux-gnu/mecab/dic/mecab-ipadic-neologd/" > /etc/mecabrc \
  && apt-get autoremove -y

ENV CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -lmecab -lstdc++"
ENV CGO_CFLAGS="-I/usr/include"

WORKDIR /app
COPY . ./
RUN go build -o /app/haiken
CMD [ "/app/haiken" ]
