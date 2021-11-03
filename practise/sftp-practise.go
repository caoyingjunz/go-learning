package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SftpConnect(user, passwd, host string, port int) (*sftp.Client, error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(passwd))

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh dial failed %v", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("new sftp client failed %v", err)
	}

	return sftpClient, nil
}

func CopyFromRemote(remoteFile, localFile, user, passwd, host string, port int) error {
	sftpClient, err := SftpConnect(user, passwd, host, port)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("open remote file %s failed %v", remoteFile, err)
	}
	defer srcFile.Close()

	f, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("create file %s failed %v", localFile, err)
	}
	defer f.Close()

	if _, err = srcFile.WriteTo(f); err != nil {
		return fmt.Errorf("write to dest file %s failed %v", localFile, err)
	}

	return nil
}

func main() {
	remoteFile := "/root/test.txt"
	localFile := "/Users/caoyuan/test.txt"

	if err := CopyFromRemote(remoteFile, localFile, "root", "passwd", "host", 22); err != nil {
		fmt.Println("copy from remote file failed", err)
	}
}
