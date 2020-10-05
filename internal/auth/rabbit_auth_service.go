package auth

import (
	"context"
	"strings"
	"time"

	"github.com/malczuuu/rmqbe/internal/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const timeout = 2 * time.Second

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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := a.database.Collection(a.collectionName).FindOne(ctx, query).Decode(&entity)
	if err != nil {
		message := "An error while retrieving user by username"
		a.logMongoFailure(err, query, message)
		result = false
	} else if entity["password"] == password {
		result = true
	}

	log.WithFields(
		log.Fields{
			"username":   username,
			"successful": result,
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := a.database.Collection(a.collectionName).FindOne(ctx, query).Decode(&entity)
	if err != nil {
		message := "An error while retrieving user by username and vhost"
		a.logMongoFailure(err, query, message)
		result = false
	}

	log.WithFields(
		log.Fields{
			"username":   username,
			"vhost":      vhost,
			"ip":         ip,
			"successful": result,
		}).Info("Authorize user to virtual host")
	return result
}

func matchNameByPattern(ref string, value string) bool {
	checkByPrefix := false
	checkBySuffix := false
	if strings.HasSuffix(ref, "*") {
		ref = strings.TrimSuffix(ref, "*")
		checkByPrefix = true
	}
	if strings.HasPrefix(ref, "*") {
		ref = strings.TrimPrefix(ref, "*")
		checkBySuffix = true
	}
	result := true

	if checkByPrefix && checkBySuffix {
		result = result && strings.Contains(value, ref)
	} else if checkByPrefix {
		result = result && strings.HasPrefix(value, ref)
	} else if checkBySuffix {
		result = result && strings.HasSuffix(value, ref)
	} else {
		result = result && value == ref
	}
	return result
}

// Resource checks whether requested user has appropriate permission to requested resource.
func (a *RabbitAuthService) Resource(username string, vhost string, resource string, name string, permission string) bool {
	result := true

	entity := bson.M{}
	query := bson.M{
		"username": username,
		"vhosts":   vhost,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := a.database.Collection(a.collectionName).FindOne(ctx, query).Decode(&entity)
	if err == nil {
		result = checkResourcePermission(entity, resource, name, permission)
	} else {
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
			"successful": result,
		}).Info("Authorize user to resource")
	return result
}

func checkResourcePermission(entity bson.M, resource string, name string, permission string) bool {
	if entity["permissions"] != nil {
		for _, v := range entity["permissions"].(bson.A) {
			value := v.(bson.M)
			if value["resource"] == resource &&
				matchNameByPattern(value["name"].(string), name) &&
				value["permission"] == permission {
				return true
			}
		}
	}
	return false
}

// Topic checks whether requested user has appropriate permission to requested topic.
func (a *RabbitAuthService) Topic(username string, vhost string, resource string, name string, permission string, routingKey string) bool {
	result := true

	entity := bson.M{}
	query := bson.M{
		"username": username,
		"vhosts":   vhost,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := a.database.Collection(a.collectionName).FindOne(ctx, query).Decode(&entity)
	if err == nil {
		result = checkTopicPermission(entity, resource, name, permission, routingKey)
	} else {
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
			"successful":  result,
		}).Info("Authorize user to topic")
	return result
}

func checkTopicPermission(entity bson.M, resource string, name string, permission string, routingKey string) bool {
	if entity["permissions"] != nil {
		for _, v := range entity["permissions"].(bson.A) {
			value := v.(bson.M)
			if value["resource"] == resource &&
				matchNameByPattern(value["name"].(string), name) &&
				value["permission"] == permission &&
				value["routing_key"] == routingKey {
				return true
			}
		}
	}
	return false
}
