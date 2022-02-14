package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var (
	errInvalidCredential = errors.New("invalid credential")
)

type MqttAuth struct {
	userStore map[string]string
}

func newMqttAuth() (*MqttAuth, error) {
	ma := &MqttAuth{
		make(map[string]string),
	}

	err := ma.readUserDb()
	return ma, err
}

func (a *MqttAuth) Authenticate(user, password []byte) bool {
	log.Infof("authenticate user: %s, %s", user, string(password))
	pwd, ok := a.userStore[string(user)]
	if !ok {
		return false
	}

	return pwd == string(password)
}

func (a *MqttAuth) ACL(user []byte, topic string, write bool) bool {
	return true
}

func (a *MqttAuth) readUserDb() error {
	file, err := os.Open("/etc/kesmarki/users")
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		u := strings.Split(scanner.Text(), ":")
		if len(u) != 2 {
			return errInvalidCredential
		}
		a.userStore[u[0]] = u[1]
	}

	return nil
}
