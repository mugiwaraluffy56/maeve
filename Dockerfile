FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/maeve ./cmd/maeve

FROM alpine:3.22
COPY --from=build /out/maeve /usr/local/bin/maeve
ENTRYPOINT ["maeve"]

