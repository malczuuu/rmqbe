package auth

import (
	"context"
	"strings"
	"time"

	"github.com/malczuuu/rmqbe/internal/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const timeout = 2 * time.Second

func NewRabbitAuthService(client *mongo.Client, config config.Config) RabbitAuthService {
	database := client.Database("rmqbe")
	return RabbitAuthService{client: client, database: database, collectionName: config.MongoUsersCollection}
}

type RabbitAuthService struct {
	client         *mongo.Client
	database       *mongo.Database
	collectionName string
}

func (a *RabbitAuthService) User(username string, password string) bool {
	result := false

	entity := bson.M{}
	query := bson.M{"username": username}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := a.database.Collection(a.collectionName).FindOne(ctx, query).Decode(&entity)
	if err != nil {
		a.logMongoFailure(err, query, "an error while retrieving user by username")
		result = false
	} else if entity["password"] == password {
		result = true
	}

	log.Info().Str("username", username).Bool("successful", result).Msg("authenticate user by username and password")
	return result
}

func (a *RabbitAuthService) logMongoFailure(err error, query bson.M, message string) {
	if err == mongo.ErrNoDocuments {
		log.Debug().Interface("query", query).Str("collection", a.collectionName).Err(err).Msg(message)
	} else {
		log.Error().Interface("query", query).Str("collection", a.collectionName).Err(err).Msg(message)
	}
}

func (a *RabbitAuthService) Vhost(username string, vhost string, ip string) bool {
	result := true

	entity := bson.M{}
	query := bson.M{"username": username, "vhosts": vhost}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := a.database.Collection(a.collectionName).FindOne(ctx, query).Decode(&entity)
	if err != nil {
		a.logMongoFailure(err, query, "an error while retrieving user by username and vhost")
		result = false
	}

	log.Info().
		Str("username", username).
		Str("vhost", vhost).
		Str("ip", ip).
		Bool("successful", result).
		Msg("authorize user to virtual host")
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
		a.logMongoFailure(err, query, "an error while retrieving user by username, vhost and resource permission")
		result = false
	}

	log.Info().
		Str("username", username).
		Str("vhost", vhost).
		Str("resource", resource).
		Str("name", name).
		Str("permission", permission).
		Bool("successful", result).
		Msg("authorize user to resource")
	return result
}

func dToM(d bson.D) bson.M {
	m := bson.M{}
	for _, elem := range d {
		m[elem.Key] = elem.Value
	}
	return m
}

func checkResourcePermission(entity bson.M, resource, name, permission string) bool {
	perms, ok := entity["permissions"].(bson.A)
	if !ok {
		return false
	}

	for _, v := range perms {
		doc, ok := v.(bson.D)
		if !ok {
			continue
		}

		value := dToM(doc)

		if value["resource"] == resource &&
			matchNameByPattern(value["name"].(string), name) &&
			value["permission"] == permission {
			return true
		}
	}

	return false
}

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

	log.Info().
		Str("username", username).
		Str("vhost", vhost).
		Str("resource", resource).
		Str("name", name).
		Str("permission", permission).
		Str("routing_key", routingKey).
		Bool("successful", result).
		Msg("authorize user to topic")
	return result
}

func checkTopicPermission(entity bson.M, resource, name, permission, routingKey string) bool {
	perms, ok := entity["permissions"].(bson.A)
	if !ok {
		return false
	}

	for _, v := range perms {
		doc, ok := v.(bson.D)
		if !ok {
			continue
		}
		value := dToM(doc)

		if value["resource"] == resource &&
			matchNameByPattern(value["name"].(string), name) &&
			value["permission"] == permission &&
			value["routing_key"] == routingKey {
			return true
		}
	}

	return false
}
