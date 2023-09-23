test:
	go test ./... -v
run:
	go run .

.PHONY:
	test run
