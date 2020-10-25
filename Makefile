all: ipquery

ipquery:
	@cd cmd/ipquery && go build -o ../../bin/ipquery

clean:
	@rm -rf bin
