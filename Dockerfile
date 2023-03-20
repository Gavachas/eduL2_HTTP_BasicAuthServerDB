FROM golang:latest
RUN mkdir app
WORKDIR /app
COPY . .
ENV EnvDBtype="sqlite"
RUN go build  ./cmd/itilium 
CMD ["./itilium"]