FROM golang:1.20.6-alpine3.18 AS build
WORKDIR /app
# Copy dependencies list
COPY go.mod go.sum ./
# Build
COPY main.go .
COPY src src
RUN go build -o main main.go

# Copy artifacts to a clean image
FROM alpine:3.18.2 AS main

COPY scripts/db-init.sh .
RUN chmod +x db-init.sh && ./db-init.sh

COPY --from=build /app/main /main
ENTRYPOINT [ "/main" ]