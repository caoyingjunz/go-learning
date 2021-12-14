package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"k8s.io/client-go/util/homedir"
)

func SftpConnect(user, passwd, host string, port int) (*sftp.Client, error) {
	// 1. 使用密码
	//clientConfig := &ssh.ClientConfig{
	//	User: user,
	//	Auth: []ssh.AuthMethod{
	//		ssh.Password(passwd),
	//	},
	//	Timeout: 30 * time.Second,
	//	HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
	//		return nil
	//	},
	//}

	// 2. 使用公钥
	key, err := ioutil.ReadFile(path.Join(homedir.HomeDir(), ".ssh", "id_rsa"))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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
