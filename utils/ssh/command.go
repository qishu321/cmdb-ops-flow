package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

type SSHClientConfig struct {
	UserName  string `form:"username" binding:"required"`
	Password  string `form:"password"`
	IP        string `form:"ip" binding:"required"`
	Port      int    `form:"port" binding:"required"`
	Command   string `form:"command"`
	AuthModel string `form:"authmodel" binding:"required"`
	PublicKey string
	Filename  string `json:"filename"`
	Content   string `json:"content"`

	Timeout time.Duration
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
func CreateFileOnRemoteServer(sshConfig *SSHClientConfig, filename, content string) (int, error) {
	absoluteFilePath := "/tmp/" + filename

	// Escape special characters and format the content for bash script
	escapedContent := strings.ReplaceAll(content, "'", `'\''`)
	formattedContent := fmt.Sprintf("#!/bin/bash\n%s", escapedContent)

	command := fmt.Sprintf("echo '%s' > %s", formattedContent, absoluteFilePath)

	output, err := SshCommand(sshConfig, command)
	if err != nil {
		fmt.Println("SSH Error:", err)
		return 400, fmt.Errorf("Failed to execute SSH command: %v", err)
	}

	fmt.Println("SSH Output:", output) // Optional: Print the SSH output

	return 200, nil
}

//func CreateFileOnRemoteServer(sshConfig *SSHClientConfig, filename, content string) (int,error) {
//	absoluteFilePath := "/tmp/" + filename
//
//	command := fmt.Sprintf("echo %q > %s", content, absoluteFilePath)
//
//	output, err := SshCommand(sshConfig, command)
//	if err != nil {
//		fmt.Println("SSH Error:", err)
//		return 400,fmt.Errorf("Failed to execute SSH command: %v", err)
//	}
//
//
//	fmt.Println("SSH Output:", output) // Optional: Print the SSH output
//
//	return 200, nil
//}
