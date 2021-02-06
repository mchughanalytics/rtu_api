package rtuapi

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// RmClient manages the connection to the reMarkable tablet device
type RmClient struct {
	Connected        bool
	Username         string
	Password         string
	Host             string
	Type             string //wifi or usb
	Connection       *ssh.Client
	ConnectionConfig *ssh.ClientConfig
}

// NewRmClient returns a new RmClient object
func NewRmClient(h, u, p string) (*RmClient, error) {

	if len(strings.Split(h, ":")) < 2 {
		h = fmt.Sprintf("%s:22", h)
	}

	rmc := &RmClient{
		Username:  u,
		Password:  p,
		Host:      h,
		Connected: false,
	}

	if h == "10.11.99.1:22" {
		rmc.Type = "usb"
	} else {
		rmc.Type = "wifi"
	}

	err := rmc.configureSSH()
	if err != nil {
		return nil, err
	}

	err = rmc.Dial()
	if err != nil {
		return nil, err
	}

	rmc.Connected = true
	//fmt.Println(rmc)

	res, err := rmc.RunCommand("echo $PWD")
	if err != nil {
		fmt.Println(err)
	} else {
		for i, line := range res {
			fmt.Printf("%d: %s", i, line)
		}

	}

	return rmc, nil
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

	rmc.ConnectionConfig = &conf

	return nil
}

func (rmc *RmClient) Dial() error {
	conn, err := ssh.Dial("tcp", rmc.Host, rmc.ConnectionConfig)
	if err != nil {
		return err
	}
	rmc.Connection = conn
	return nil
}

//RunCommand against ssh connection
func (rmc *RmClient) RunCommand(cmd string) ([]string, error) {

	log.Println("Running command")

	output := new(strings.Builder)

	sess, err := rmc.Connection.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return nil, err
	}

	//go io.Copy(os.Stdout, sessStdOut)
	go io.Copy(output, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return nil, err
	}

	//go io.Copy(os.Stderr, sessStderr)
	go io.Copy(output, sessStderr)

	//log.Printf("Executing: %s", cmd)
	err = sess.Run(cmd)
	time.Sleep(1 * time.Second)
	if err != nil {
		return nil, err
	}

	outRaw := strings.Split(output.String(), "\n")
	out := []string{}

	for _, line := range outRaw {
		if len(line) > 0 {
			out = append(out, line)
		}
	}

	return out, nil
}
