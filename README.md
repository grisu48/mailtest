# mailtest

Simple go test program for mailserver

I wrote this program to test the configuration and login of my own mailservers

## Requirements

    go get github.com/emersion/go-imap/...

## Usage

    go build mailtest.go
    ./mailtest mailserver.com:993

The program tests a connection and lets you type in your credentials. In the best case it prints the mailboxes or it spits out an error, if something went wrong
