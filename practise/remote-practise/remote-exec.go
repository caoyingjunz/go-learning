package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"golang.org/x/crypto/ssh"
	"k8s.io/client-go/util/homedir"
)

func RemoteConnect(user, passwd, host string, port int) (*ssh.Session, error) {
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

	session, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func main() {
	var stdOut, stdErr bytes.Buffer

	session, err := RemoteConnect("root", "passwd", "peng", 22)
	if err != nil {
		fmt.Println("copy from remote file failed", err)
	}
	defer session.Close()

	session.Stdout = &stdOut
	session.Stderr = &stdErr

	if err = session.Run("ls -al"); err != nil {
		fmt.Println(stdErr.String())
		return
	}

	fmt.Println(stdOut.String())
}
