build:
	go build -o gocurl main.go
test:
	./gocurl -v http://eu.httpbin.org/get