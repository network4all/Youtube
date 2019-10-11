package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

const (
	username = "admin"
	password = "c"
	hostname = "192.168.0.200"
	port     = ":22"
	command  = "send log 6 \"Hello world from google go!\""
	//command = "show run"
)

func main() {

	// create command
	cmds := []string{"terminal length 0",
		command,
		"exit",
	}

	// config
	config := &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: []string{"aes128-cbc"},
		},
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Printf("Connecting to device %s%s... \n", hostname, port)

	// make tcp connection
	device, err := ssh.Dial("tcp", hostname+port, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer device.Close()

	// create session
	fmt.Printf("Creating session... \n")
	session, err := device.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	//redirect pipes
	stdin, _ := session.StdinPipe()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Shell()

	// send commands
	for _, cmd := range cmds {
		fmt.Printf("Sending command: '%s'\n", cmd)
		fmt.Fprintf(stdin, "%s\n", cmd)
	}

	// wait for disconnect
	fmt.Printf("Waiting for session disconnect...\n")
	session.Wait()
}
