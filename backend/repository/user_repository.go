package repository

import "sync"

var (
	users = make(map[string]string)
	mutex sync.Mutex
)

func InsertUser(name, password string) error {
	mutex.Lock()
	defer mutex.Unlock()
	users[name] = password
	return nil
}

func GetUsers() map[string]string {
	return users
}

func RemoveUser(name string) error {
	mutex.Lock()
	defer mutex.Unlock()
	delete(users, name)
	return nil
}
