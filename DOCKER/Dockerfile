# docker build --no-cache -t node:latest -f DOCKER/Dockerfile .
# docker build -t noah-extender:latest -f DOCKER/Dockerfile .
# docker run -d -p 127.0.0.1:9000:9000 --restart=always noah-extender:latest

FROM golang:1.12-buster as builder

ENV APP_PATH /home/CoinExplorer-BackEnd

COPY . ${APP_PATH}

WORKDIR ${APP_PATH}

RUN make create_vendor && \
    make build && \
    cp ./builds/coin-explorer /usr/local/bin/coin-explorer

EXPOSE 9070
CMD ["coin-explorer"]
STOPSIGNAL SIGTERM
