all:gen

gen:
	@thrift -gen   go thrift/servant.thrift
	@thrift -gen   py thrift/servant.thrift
	@thrift -gen  php thrift/servant.thrift
	@thrift -gen java thrift/servant.thrift

clean:
	@rm -rf gen-go
	@rm -rf gen-php
	@rm -rf gen-py
	@rm -rf gen-java
