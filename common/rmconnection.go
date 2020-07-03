package common

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

// RmClient manages the connection to the reMarkable tablet device
type RmClient struct {
	Connected  bool
	Username   string
	Password   string
	Host       string
	Type       string //wifi or usb
	Connection ssh.Client
}

// NewRmClient returns a new RmClient object
func NewRmClient(host, username, password string) (*RmClient, error) {
	rmc := RmClient{
		Username: username,
		Password: password,
		Host:     host,
	}

	if host == "10.11.99.1" {
		rmc.Type = "usb"
	} else {
		rmc.Type = "wifi"
	}

	err := rmc.configureSSH()
	if err != nil {
		return nil, err
	}

	return &rmc, nil
}

func (rmc *RmClient) configureSSH() error {

	a := []ssh.AuthMethod{}
	a = append(a, ssh.Password(rmc.Password))

	conf := ssh.ClientConfig{
		User:            rmc.Username,
		Auth:            a,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         0,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", rmc.Host), &conf)
	if err != nil {
		return err
	}
	defer conn.Close()
	rmc.Connection.Conn = conn

	err = rmc.RunCommand("ls -la /usr/share/remarkable/templates/")
	if err != nil {
		return err
	}

	return nil
}

//RunCommand against ssh connection
func (rmc *RmClient) RunCommand(cmd string) error {

	sess, err := rmc.Connection.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return err
	}

	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return err
	}

	go io.Copy(os.Stderr, sessStderr)
	err = sess.Run(cmd) // eg., /usr/bin/whoami
	if err != nil {
		return err
	}

	return nil
}
