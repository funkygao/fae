all:gen

gen:
	thrift -gen   go servant.thrift
	thrift -gen   py servant.thrift
	thrift -gen  php servant.thrift
	thrift -gen java servant.thrift

clean:
	rm -rf gen-go
	rm -rf gen-php
	rm -rf gen-py
