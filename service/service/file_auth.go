package service

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var (
	errInvalidCredential = errors.New("invalid credential")
)

type FileAuth struct {
	userStore map[string]string
}

func NewFileAuth(dbFile string) (*FileAuth, error) {
	ma := &FileAuth{
		make(map[string]string),
	}

	err := ma.readUserDb(dbFile)
	return ma, err
}

func (a *FileAuth) Authenticate(user, password []byte) bool {
	log.Infof("authenticate user: %s, %s", user, string(password))
	pwd, ok := a.userStore[string(user)]
	if !ok {
		return false
	}

	return pwd == string(password)
}

func (a *FileAuth) ACL(user []byte, topic string, write bool) bool {
	return true
}

func (a *FileAuth) readUserDb(dbFile string) error {
	file, err := os.Open(dbFile)
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
