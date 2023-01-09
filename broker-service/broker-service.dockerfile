# base go image
# FROM golang:1.18-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

# # build a tiny docker image
# FROM alpine:latest

# RUN mkdir /app

# COPY --from=builder /app/brokerApp /app

# CMD [ "/app/brokerApp" ]


# Since, we are using Makefile, we don't need the above dockerfile configuration anymore
# Because, we are building linux executable of our go application while using
# "make up_build" or "make build_broker" command. So, we only need to copy this
# executable file(brokerApp) to the /app workdir of our container(alpine image). No need to build
# multi stage dockerfile and this will take much less time than the multi stage build
FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD [ "/app/brokerApp" ]

