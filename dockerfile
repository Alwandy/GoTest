FROM golang:latest
WORKDIR /
COPY . . 
RUN go get -d github.com/gorilla/mux 
RUN go get -d github.com/go-sql-driver/mysql

CMD ["go", "run", "main.go"]