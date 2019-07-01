default: all
all: mailtest

mailtest: mailtest.go
	go build mailtest.go
install: mailtest
	install mailtest /usr/local/bin/mailtest
