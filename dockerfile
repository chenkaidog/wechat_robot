FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go env -w GOPROXY='https://goproxy.io/' && \ 
    go mod tidy && \
    go build -o main && \
    mkdir log storage

ENV log_output_file_name=openwechat \
    log_set_local_file=true \
    log_level=debug

CMD [ "./main" ]