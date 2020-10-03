package auth

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAuthManager(client *mongo.Client) AuthManager {
	return AuthManager{client: client}
}

type AuthManager struct {
	client *mongo.Client
}

func (a *AuthManager) User(username string, password string) bool {
	result := false

	if username != "" && password != "" {
		result = true
	}

	log.WithFields(
		map[string]interface{}{
			"username": username,
			"result":   result,
		}).Info("Authenticate user by username and password")
	return result
}

func (a *AuthManager) Vhost(username string, vhost string, ip string) bool {
	result := true

	if vhost == "/" && username != "" {
		result = false
	}

	log.WithFields(
		map[string]interface{}{
			"username": username,
			"vhost":    vhost,
			"ip":       ip,
			"result":   result,
		}).Info("Authorize user to virtual host")
	return result
}

func (a *AuthManager) Resource(username string, vhost string, resource string, name string, permission string) bool {
	result := true
	log.WithFields(
		map[string]interface{}{
			"username":   username,
			"vhost":      vhost,
			"resource":   resource,
			"name":       name,
			"permission": permission,
			"result":     result,
		}).Info("Authorize user to resource")
	return result
}

func (a *AuthManager) Topic(username string, vhost string, resource string, name string, permission string, routingKey string) bool {
	result := true
	log.WithFields(
		map[string]interface{}{
			"username":    username,
			"vhost":       vhost,
			"resource":    resource,
			"name":        name,
			"permission":  permission,
			"routing_key": routingKey,
			"result":      result,
		}).Info("Authorize user to topic")
	return result
}
