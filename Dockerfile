FROM golang:1.20.7-alpine

COPY . project
RUN cd project && go build -o http-proxy http-proxy/app

EXPOSE 8080/tcp
EXPOSE 80/tcp

ENTRYPOINT [ "./project/http-proxy" ]
