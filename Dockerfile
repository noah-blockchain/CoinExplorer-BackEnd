# docker build --no-cache -t noah-explorer:latest -f Dockerfile .
# docker build -t noah-explorer:latest -f Dockerfile .

FROM golang:1.12-buster as builder
ENV APP_PATH /home/noah-explorer-api
COPY . ${APP_PATH}
WORKDIR ${APP_PATH}
RUN make create_vendor && make build

FROM debian:buster-slim as executor
COPY --from=builder /home/noah-explorer-api/build/coin-explorer /usr/local/bin/coin-explorer
CMD ["coin-explorer"]
STOPSIGNAL SIGTERM
