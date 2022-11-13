package mqtt

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var (
	errInvalidCredential = errors.New("invalid credential")
	errUsersDbNotFound   = errors.New("user store file not found")
)

type FileAuth struct {
	userStore map[string]string
}

func NewFileAuth(dbFiles ...string) (*FileAuth, error) {
	ma := &FileAuth{
		make(map[string]string),
	}

	f, ok := ma.chooseDbFile(dbFiles)
	if !ok {
		return ma, errUsersDbNotFound
	}

	err := ma.readUserDb(f)
	return ma, err
}

func (a *FileAuth) Authenticate(user, password []byte) bool {
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

func (a *FileAuth) chooseDbFile(dbFiles []string) (string, bool) {
	for _, s := range dbFiles {
		if _, err := os.Stat(s); err == nil {
			return s, true
		}
	}

	return "", false
}
