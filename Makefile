# Usage:
# make			# make build
# make build	# build binary
# make clean	# delete binary directory
# make run		# run code

BIN_DIR=bin/

build: clean
	@go build -v -a -o $(BIN_DIR) -race

clean:
	@rm -rf $(BIN_DIR)

run:
	@go run main.go