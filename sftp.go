package rtuapi

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RmSFTP struct {
	client        *sftp.Client
	clientOptions []*sftp.ClientOption
}

func (ftp *RmSFTP) ReadRemoteFile(path string) []byte {
	return nil
}

func NewRmSFTP(ssh *ssh.Client) *RmSFTP {
	return nil
}
