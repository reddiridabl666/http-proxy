FROM golang:1.20.7-alpine

COPY . project

RUN cd project && \
    go build -o http-proxy http-proxy/app && \
    chmod +x gen.sh && \
    mkdir certs && \
    apk add openssl

EXPOSE 8080/tcp
EXPOSE 8000/tcp

CMD cd project && ./http-proxy
