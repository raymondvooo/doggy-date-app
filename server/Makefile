run: bin/doggy-date-go
	@PATH="$(PWD)/bin:$(PATH)" heroku local

bin/doggy-date-go: server.go
	go build -o bin/bin/doggy-date-go server.go

clean:
	rm -rf bin