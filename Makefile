test:
	go test ./... -v
run:
	go run .
install:
	go install .

.PHONY:
	test run
