package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TableName = "artists"

func (d *DynamoRepository) GetArtist(ctx context.Context, id uuid.UUID) (*setmakerpb.Artist, error) {
	// create key map
	keys, err := attributevalue.MarshalMap(map[string]string{
		"Id": *aws.String(id.String()),
	})
	if err != nil {
		logger.WithField("id", id).Errorf("GetArtist: Could not marshalmap: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid UUID")
	}

	// fetch item from dynamo
	data, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       keys,
	})
	if err != nil {
		logger.WithField("id", id).Errorf("GetArtist: Error fetching from dynamo: %s", err)
		return nil, status.Error(codes.Internal, "Error fetching result")
	}

	// check an item was returned
	if data.Item == nil {
		logger.WithField("id", id).Error("GetArtist: No artist found for ID")
		return nil, status.Error(codes.NotFound, "Artist not found")
	}

	// fetch was successful
	logger.WithField("data", data.Item).Info("Artist found")

	// unmarshal response
	var res *setmakerpb.Artist
	err = attributevalue.UnmarshalMap(data.Item, res)
	if err != nil {
		logger.WithField("data", data.Item).Errorf("GetArtist: Could not unmarshal item: %s", err)
		return nil, fmt.Errorf("GetArtist: Error unmarshaling artist data")
	}

	return res, nil
}

func (d *DynamoRepository) PutArtist(ctx context.Context, artist *setmakerpb.Artist) error {
	// create attribute value map
	item, err := attributevalue.MarshalMap(artist)
	if err != nil {
		logger.WithField("data", artist).Errorf("GetArtist: Could not marshalmap: %s", err)
		return fmt.Errorf("GetArtist: Could not map input values for artist")
	}

	// PutItem to dynamo
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})
	if err != nil {
		logger.WithField("data", artist).Errorf("Could not PutItem: %s", err)
		return fmt.Errorf("GetArtist: Failed to persist artist")
	}

	logger.WithField("id", artist.Id).Info("Artist persisted successfully")
	return nil
}
