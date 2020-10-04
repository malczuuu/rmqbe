package auth

import (
	"strings"

	"github.com/malczuuu/rmqbe/internal/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewRabbitAuthService creates new instance of RabbitAuthService.
func NewRabbitAuthService(client *mongo.Client, config config.Config) RabbitAuthService {
	database := client.Database("rmqbe")
	return RabbitAuthService{client: client, database: database, collectionName: config.MongoUsersCollection}
}

// RabbitAuthService instance is the main engine of RabbitMQ authorization service.
type RabbitAuthService struct {
	client         *mongo.Client
	database       *mongo.Database
	collectionName string
}

// User checks whether requested user exists and it's password is correct.
func (a *RabbitAuthService) User(username string, password string) bool {
	result := false

	entity := bson.M{}
	query := bson.M{"username": username}
	err := a.database.Collection(a.collectionName).FindOne(nil, query).Decode(&entity)
	if err != nil {
		message := "An error while retrieving user by username"
		a.logMongoFailure(err, query, message)
		result = false
	} else if entity["password"] == password {
		result = true
	}

	log.WithFields(
		log.Fields{
			"username": username,
			"result":   result,
		}).Info("Authenticate user by username and password")
	return result
}

func (a *RabbitAuthService) logMongoFailure(err error, query bson.M, message string) {
	preparedLog := log.WithFields(log.Fields{"query": query, "collection": a.collectionName}).WithError(err)
	if err == mongo.ErrNoDocuments {
		preparedLog.Debug(message)
	} else {
		preparedLog.Error(message)
	}
}

// Vhost checks whether requested user is a member of desired virtual host.
func (a *RabbitAuthService) Vhost(username string, vhost string, ip string) bool {
	result := true

	entity := bson.M{}
	query := bson.M{"username": username, "vhosts": vhost}
	err := a.database.Collection(a.collectionName).FindOne(nil, query).Decode(&entity)
	if err != nil {
		message := "An error while retrieving user by username and vhost"
		a.logMongoFailure(err, query, message)
		result = false
	}

	log.WithFields(
		log.Fields{
			"username": username,
			"vhost":    vhost,
			"ip":       ip,
			"result":   result,
		}).Info("Authorize user to virtual host")
	return result
}

// Resource checks whether requested user has appropriate permission to requested resource.
func (a *RabbitAuthService) Resource(username string, vhost string, resource string, name string, permission string) bool {
	// TODO: add pattern support and check for mqtt-subscription-* queue permission
	if resource == "queue" && strings.HasPrefix(name, "mqtt-subscription-") {
		return true
	}

	result := true

	entity := bson.M{}
	query := bson.M{
		"username":       username,
		"vhosts":         vhost,
		"perms.resource": resource,
		"perms.name":     name,
		"perms.perm":     permission,
	}
	err := a.database.Collection(a.collectionName).FindOne(nil, query).Decode(&entity)
	if err != nil {
		message := "An error while retrieving user by username, vhost and resource permission"
		a.logMongoFailure(err, query, message)
		result = false
	}

	log.WithFields(
		log.Fields{
			"username":   username,
			"vhost":      vhost,
			"resource":   resource,
			"name":       name,
			"permission": permission,
			"result":     result,
		}).Info("Authorize user to resource")
	return result
}

// Topic checks whether requested user has appropriate permission to requested topic.
func (a *RabbitAuthService) Topic(username string, vhost string, resource string, name string, permission string, routingKey string) bool {
	result := true

	entity := bson.M{}
	query := bson.M{
		"username":          username,
		"vhosts":            vhost,
		"perms.resource":    resource,
		"perms.name":        name,
		"perms.perm":        permission,
		"perms.routing_key": routingKey,
	}
	err := a.database.Collection(a.collectionName).FindOne(nil, query).Decode(&entity)
	if err != nil {
		message := "An error while retrieving user by username, vhost and topic permission"
		a.logMongoFailure(err, query, message)
		result = false
	}

	log.WithFields(
		log.Fields{
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
