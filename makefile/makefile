BINARY_NAME=PackagebirdServer.exe
BUILD_PATH='build/$(BINARY_NAME)'


build:
	go build -o ${BUILD_PATH} main.go
run:
	go build -o ${BUILD_PATH} main.go
	./${BUILD_PATH}
clean:
	go clean
	rm ${BUILD_PATH}
	rm 'build/builds'
	rm 'build/packages'
	rm 'build/projects'
	rm 'build/tmp'
	$(shell mongo packagebird --eval "db.dropDatabase()")

