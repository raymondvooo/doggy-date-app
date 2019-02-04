run: bin/server
	@PATH="$(PWD)/bin/doggy-date-go" heroku local

bin/server: server.go
	go build -v -o bin/server server.go

clean:
	rm -rf bin