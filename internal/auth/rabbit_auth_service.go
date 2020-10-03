package auth

import (
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
	err := a.database.Collection("users").FindOne(nil, query).Decode(&entity)
	if err != nil {
		log.WithFields(log.Fields{
			"username": username,
			"query":    query,
		}).WithError(err).Debug("An error while retrieving user by username")
		if err != mongo.ErrNoDocuments {
			log.WithField("username", username).WithError(err).Error("Failed to authenticate user by password")
		}
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

// Vhost checks whether requested user is a member of desired virtual host.
func (a *RabbitAuthService) Vhost(username string, vhost string, ip string) bool {
	result := true

	entity := bson.M{}
	query := bson.M{"username": username, "vhosts": vhost}
	err := a.database.Collection("users").FindOne(nil, query).Decode(&entity)
	if err != nil {
		log.WithField("query", query).WithError(err).Debug("An error while retrieving user by username and vhost")
		if err != mongo.ErrNoDocuments {
			log.WithError(err).Error("Failed to authorize user to virtual host")
		}
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
	result := true

	entity := bson.M{}
	query := bson.M{
		"username": username,
		"vhosts":   vhost,
		"permissions": bson.M{
			"resource":   resource,
			"name":       name,
			"permission": permission,
		},
	}
	err := a.database.Collection("users").FindOne(nil, query).Decode(&entity)
	if err != nil {
		log.WithField("query", query).WithError(err).Debug("An error while retrieving user by username, vhost and resource permission")
		if err != mongo.ErrNoDocuments {
			log.WithError(err).Error("Failed to authorize user to resource")
		}
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
		"username": username,
		"vhosts":   vhost,
		"permissions": bson.M{
			"resource":    resource,
			"name":        name,
			"permission":  permission,
			"routing_key": routingKey,
		},
	}
	err := a.database.Collection("users").FindOne(nil, query).Decode(&entity)
	if err != nil {
		log.WithField("query", query).WithError(err).Debug("An error while retrieving user by username, vhost and topic permission")
		if err != mongo.ErrNoDocuments {
			log.WithError(err).Error("Failed to authorize user to topic")
		}
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
