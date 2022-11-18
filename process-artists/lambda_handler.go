package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/golang/protobuf/proto"
	"github.com/pete-robinson/setmaker-artist-meta-lambda/process-artists/repository"
	"github.com/pete-robinson/setmaker-artist-meta-lambda/process-artists/service"
	"github.com/pete-robinson/setmaker-artist-meta-lambda/process-artists/utils"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	EnvAwsAccessKey        = "AWSACCESS_KEY_ID"
	EnvAwsAccessSecret     = "AWSSECRET_ACCESS_KEY"
	EnvAwsRegion           = "AWSREGION"
	EnvSpotifyClientId     = "SPOTIFY_CLIENT_ID"
	EnvSpotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
)

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		msgStr := snsRecord.Message

		logger.WithFields(logger.Fields{
			"eventSource": record.EventSource,
			"timestamp": snsRecord.Timestamp,
			"message": snsRecord.Message,
		}).Info("Message Received. Beginning processing")

		msg := []byte(msgStr) 		// str to byte slice
		bootstrap(ctx, msg)				// kick off the lambda handler
	}
}

func main() {
	lambda.Start(handler)
}

// @todo - refactor service creation to outside of the handler record loop - not efficient to build a new service struct for every event in the payload
func bootstrap(ctx context.Context, msg []byte) {
	logger.WithField("msg", string(msg)).Info("Beginning processing for message")

	// unmarshal the message struct
	evt := &setmakerpb.Event{}

	err := proto.Unmarshal(msg, evt)
	if err != nil {
		logger.WithField("msg", string(msg)).Fatalf("bootstrap: Message could not be unmarshaled: %s", err)
		panic(err)
	}

	logger.WithField("event", evt).Info("Event unmarshaled successfully")

	switch evt.EventType {
		case setmakerpb.Event_EVENT_ARTIST_CREATED:
			logger.Info("Event identified as type EVENT_ARTIST_CREATED")
			handleArtistCreate(ctx, evt.GetArtistCreated())

		case setmakerpb.Event_EVENT_ARTIST_DELETED:
			logger.Info("Event identified as type EVENT_ARTIST_DELETED")
			panic("NOT IMPLEMENTED")

		default:
			break
	}
}

func handleArtistCreate(ctx context.Context, msg *setmakerpb.MessageBody_ArtistCreated) {
	// @todo - create a service bootstrap using functional options to minimise code duplication

	// init AWS and create dynamo client
	dynamoClient, err := createDynamoClient(ctx)
	if err != nil {
		panic(err)
	}

	// init repo
	repo := repository.NewDynamoRepository(dynamoClient)

	// init spotify client
	spotifyClient, err := createSpotifyClient(ctx)
	if err != nil {
		panic(err)
	}

	// fetch the relevant information from the event
	id := msg.Id
	logger.Infof("Fetching artist information for artist: %s", msg.Name)

	// init and call service
	s := service.NewService(repo, spotifyClient)
	ok, err := s.FetchArtistMeta(ctx, id)
	if err != nil {
		panic(fmt.Errorf("handleArtistCreate: %s", err))
	}

	if ok {
		logger.WithField("id", id).Infof("Artist updated successfully")
	} else {
		logger.WithField("id", id).Fatalf("handleArtistCreate: Artist was not updated")
	}
}

func createDynamoClient(ctx context.Context) (*dynamodb.Client, error) {
	awsConfigObj := &utils.AwsConfig{
		Region: os.Getenv(EnvAwsRegion),
	}

	awsConfig, err := utils.BuildAwsConfig(ctx, awsConfigObj)
	if err != nil {
		logger.Errorf("createDynamoClient: Could not build AWS config: %s", err)
		return nil, err
	}

	return utils.CreateDynamoClient(awsConfig), nil
}

func createSpotifyClient(ctx context.Context) (*service.SpotifyClient, error) {
	// init the spotify client
	config := &clientcredentials.Config{
		ClientID:     os.Getenv(EnvSpotifyClientId),
		ClientSecret: os.Getenv(EnvSpotifyClientSecret),
		TokenURL:     spotifyauth.TokenURL,
	}

	client, err := service.NewSpotifyClient(ctx, config)
	if err != nil {
		logger.Errorf("createSpotifyClient: Could not initialise Spotify client: %s", err)
		return nil, err
	}

	return client, nil
}
