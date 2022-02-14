package main

type MqttAuth struct{}

func (a *MqttAuth) Authenticate(user, password []byte) bool {
	log.Printf("user: %s, %s", user, string(password))
	return true
}

func (a *MqttAuth) ACL(user []byte, topic string, write bool) bool {
	return true
}
