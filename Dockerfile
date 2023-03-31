FROM golang:alpine AS builder
RUN apk add build-base
LABEL stage=gobuilder
ENV CGO_ENABLED=1
ENV GOOS linux
ENV ENVDBTYPE="sqlite3"
ENV ADR_GRPC="172.17.0.2:9090"
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
RUN mkdir app
RUN mkdir app/internal
RUN mkdir app/internal/database
COPY . .
COPY ./internal/database   /app/internal/database
RUN go build -o /app/itilium ./cmd/itilium 

FROM alpine
RUN mkdir app
RUN mkdir app/internal
RUN mkdir app/internal/database
EXPOSE 4000
WORKDIR /app
COPY --from=builder /app/itilium /app/itilium
COPY --from=builder  /app/internal/database   /app/internal/database
CMD ["./itilium"]