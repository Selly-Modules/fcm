package fcm

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/Selly-Modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/option"
)

// MongoDBConfig ...
type MongoDBConfig struct {
	Host, User, Password, DBName, Mechanism, Source string
}

// Config ...
type Config struct {
	// Project id
	ProjectID string
	// Original is JSON format, but encoded with base64, need to call base64 decode to get byte data
	Credential string
	// MongoDB config, for save logs
	MongoDB MongoDBConfig
}

// Result of each send
type Result struct {
	BatchID     string
	Success     int
	Failure     int
	ErrorTokens []string
}

// List topics
const (
	TopicAll     = "all"
	TopicIOS     = "iOS"
	TopicAndroid = "android"
)

var allowedTopics = []string{TopicAll, TopicIOS, TopicAndroid}

// Service ...
type Service struct {
	Config
	Client *messaging.Client
	DB     *mongo.Database
}

var s *Service

// NewInstance for push notification
func NewInstance(config Config) error {
	if config.ProjectID == "" || config.Credential == "" || config.MongoDB.Host == "" {
		return errors.New("please provide all information that needed: projectId, credential, postgresql")
	}

	ctx := context.Background()

	// Connect MongoDB
	err := mongodb.Connect(
		config.MongoDB.Host,
		config.MongoDB.User,
		config.MongoDB.Password,
		config.MongoDB.DBName,
		config.MongoDB.Mechanism,
		config.MongoDB.Source,
	)
	if err != nil {
		fmt.Println("FCM - Connect MongoDB error", err)
		return err
	}

	// Setup
	credential := base64Decode(config.Credential)
	opts := option.WithCredentialsJSON(credential)
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: config.ProjectID,
	}, opts)
	if err != nil {
		fmt.Println("FCM - New app error", err)
		return errors.New("error when init Firebase app")
	}

	// Init messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Println("FCM - init messaging client error", err)
		return errors.New("error when init Firebase messaging client")
	}

	s = &Service{
		Config: config,
		Client: client,
		DB:     mongodb.GetInstance(),
	}
	return nil
}

// GetInstance ...
func GetInstance() *Service {
	return s
}
