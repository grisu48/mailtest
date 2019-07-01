package main

import (
	"fmt"
	"os"
	"bufio"
	"syscall"
	"strings"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
	"golang.org/x/crypto/ssh/terminal"
)


func readCredentials() (string, string) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Username: ")
    username, _ := reader.ReadString('\n')

    fmt.Print("Password: ")
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil { panic(err) }
    fmt.Println()
    password := string(bytePassword)
    return strings.TrimSpace(username), strings.TrimSpace(password)
}

func main() {
	remote := ""
	useTLS := 0
	
	for _, arg := range(os.Args[1:]) {
		if arg == "--help" || arg == "-h" {
			fmt.Printf("Usage: %s [OPTIONS] REMOTE\n", os.Args[0])
			fmt.Println("OPTIONS:")
			fmt.Println("  -h, --help                   Print this help message")
			fmt.Println("      --notls, --nossl         Don't use TLS/SSL")
			fmt.Println("      --tls, --ssl             Use TLS/SSL")
			fmt.Println("REMOTE is in the form HOSTNAME:PORT, e.g. mailserver.com:143")
			os.Exit(0)
		} else if arg == "--notls" || arg == "--nossl" {
			useTLS = 1
		} else if arg == "--tls" || arg == "--ssl" {
			useTLS = 2
		} else {
			remote = arg
		}
	}
	
	if remote == "" {
		fmt.Fprint(os.Stderr, "Usage: %s REMOTE\n", os.Args[0])
		fmt.Fprint(os.Stderr, "       e.g. %s mail.example.org:993\n")
		os.Exit(1)
	}
	
	// Determine TLS, if not defined
	if useTLS == 0 {
		if strings.HasSuffix(remote, ":143") {
			useTLS = 1
		} else if strings.HasSuffix(remote, ":993") {
			useTLS = 2
		} else {
			fmt.Fprintln(os.Stderr, "Cannot determine TLS status. Please set --ssl or --nossl")
			os.Exit(1)
		}
	}

	// Connect to server
	var c *client.Client
	var err error
	if useTLS == 1 {
		c, err = client.Dial(remote)
	} else {
		c, err = client.DialTLS(remote, nil)
	}
	if err != nil {
		fmt.Fprint(os.Stderr, "Connection failed:\n", err)
		os.Exit(1)
	}
	defer c.Logout()
	fmt.Println("Connected. Please login")

	username, password := readCredentials()
	if err := c.Login(username, password); err != nil {
		fmt.Fprint(os.Stderr, "Login failed: %s\n", err)
		os.Exit(1)
	}

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func () {
		done <- c.List("", "*", mailboxes)
	}()
	fmt.Println("Mailboxes:")
	for m := range mailboxes {
		fmt.Println("\t* " + m.Name)
	}
	if err := <-done; err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching mailboxes: ", err)
		os.Exit(1)
	}

	fmt.Println("All good")
}
