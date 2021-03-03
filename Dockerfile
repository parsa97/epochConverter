FROM golang:1.15 as builder
ENV GOPATH /opt/
COPY ./epochConvertor.go ./helper.go /opt/epochConvertor/
COPY ./metrics /opt/epochConvertor/metrics
RUN cd /opt/epochConvertor \
    && go mod init epochConvertor \
    && CGO_ENABLED=0 go build epochConvertor
FROM alpine:3.12
RUN mkdir /epochConvertor
WORKDIR /epochConvertor
COPY --from=builder /opt/epochConvertor/epochConvertor /epochConvertor/
CMD ["./epochConvertor"]