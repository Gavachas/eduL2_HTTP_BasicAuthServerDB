FROM golang:latest
RUN mkdir app
WORKDIR /app
COPY . .
ENV ENVDBTYPE="sqlite"
RUN go build  ./cmd/itilium 
CMD ["./itilium"]