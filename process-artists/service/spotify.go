package service

import (
	"context"
	"fmt"

	logger "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

const ResultLimit = 5

type SpotifyClient struct {
	client *spotify.Client
}

func NewSpotifyClient(ctx context.Context, config *clientcredentials.Config) (*SpotifyClient, error) {
	// init the base client
	client, err := createClient(ctx, config)
	if err != nil {
		return nil, err
	}

	return &SpotifyClient{
		client: client,
	}, nil
}

func (s *SpotifyClient) SearchForArtist(ctx context.Context, term string) (*spotify.FullArtist, error) {
	// search API
	results, err := s.client.Search(ctx, term, spotify.SearchTypeArtist, spotify.Limit(ResultLimit))
	if err != nil {
		logger.WithField("searchTerm", term).Errorf("SearchForArtist: Search for artist failed: %s", err)
		return nil, err
	}

	// check for results and find most popular artist from resultset
	if results.Artists != nil {
		logger.WithField("searchTerm", "term").Infof("Artist found with spotify search. Found %d results", len(results.Artists.Artists))

		var res spotify.FullArtist
		mostPopular := 0

		// would be far better to use quicksort and pop the slice but...
		for _, item := range results.Artists.Artists {
			if item.Popularity > mostPopular {
				mostPopular = item.Popularity
				res = item
			}
		}

		logger.WithFields(logger.Fields{
			"searchTerm":       term,
			"MostPopularMatch": res.Name,
		}).Infof("Identified most popular artist")

		return &res, nil
	}

	logger.WithField("searchTerm", term).Errorf("No seach results found for query")
	return nil, fmt.Errorf("SearchForArtist: No search results found for query: %s", term)
}

func createClient(ctx context.Context, config *clientcredentials.Config) (*spotify.Client, error) {
	token, err := config.Token(ctx)
	if err != nil {
		logger.Errorf("createClient: Could not generate Spotify API token: %s", err)
		return nil, err
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	return spotify.New(httpClient), nil
}
