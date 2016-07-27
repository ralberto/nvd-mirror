BINARY=nvd-mirror


all: 
	go build -o ${BINARY} main.go

clean:
	@rm -f ${BINARY}

#	${PWD}/../.. go build -o ${BINARY} main.go
