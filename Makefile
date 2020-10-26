all: ipquery

ipquery:
	@cd cmd/ipquery && go build -o ../../bin/ipquery

clean:
	@rm -rf bin

install:
	mkdir -p /usr/local/ipquery/bin
	mkdir -p /var/log/ipquery
	chown nobody.www-data /usr/local/ipquery/run
	chown nobody.www-data /var/log/ipquery
	cp bin/ipquery /usr/local/ipquery/bin/
	cp systemd/ipquery.service /etc/systemd/system/ipquery.service
