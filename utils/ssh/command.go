package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHClientConfig struct {
	UserName  string `form:"username" binding:"required"`
	Password  string `form:"password"`
	IP        string `form:"ip" binding:"required"`
	Port      int    `form:"port" binding:"required"`
	Command   string `form:"command" binding:"required"`
	AuthModel string `form:"authmodel" binding:"required"`
	PublicKey string
	Timeout   time.Duration
}

func SshCommand(conf *SSHClientConfig, command string) (string, error) {
	config := &ssh.ClientConfig{
		User:            conf.UserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略know_hosts检查
	}
	switch conf.AuthModel {
	case "PASSWORD":
		config.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	case "PUBLICKEY":
		signer, err := ssh.ParsePrivateKey([]byte(conf.PublicKey))
		if err != nil {
			return "", err
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.IP, conf.Port), config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
