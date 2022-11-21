package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pete-robinson/setmaker-artist-meta-lambda/process-artists/utils"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
)

type Repository interface {
	GetArtist(context.Context, uuid.UUID) (*setmakerpb.Artist, error)
	PutArtist(context.Context, *setmakerpb.Artist) error
}

type Service struct {
	spotifyClient *SpotifyClient
	repository    Repository
}

func NewService(repo Repository, spotifyClient *SpotifyClient) *Service {
	return &Service{
		spotifyClient: spotifyClient,
		repository:    repo,
	}
}

func (s *Service) FetchArtistMeta(ctx context.Context, id string) (bool, error) {
	// parse str to UUID
	uuid, err := uuid.Parse(id)
	if err != nil {
		logger.WithField("id", id).Errorf("FetchArtistMeta: Could not convert provided ID to UUID: %s", err)
		return false, err
	}

	// fetch artist
	artist, err := s.repository.GetArtist(ctx, uuid)
	if err != nil {
		return false, err
	}

	// throw artist name at spotify API
	res, err := s.spotifyClient.SearchForArtist(ctx, artist.Name)
	if err != nil {
		return false, err
	}

	// check again for no results
	if res != nil {
		logger.WithField("searchResult", res).Info("Assigning artist data to struct")

		// update artist details
		utils.SetMetaData(artist.Metadata)
		artist.Genres = res.Genres
		artist.SpotifyUrl = string(res.URI)
		if len(res.Images) > 0 {
			artist.Image = res.Images[0].URL
		}

		logger.WithField("artist", artist).Infof("Artist data set. Attempting persistance")

		// persist artist
		if err = s.repository.PutArtist(ctx, artist); err != nil {
			return false, err
		}

		logger.WithField("artist", artist).Infof("Artist persisted without error")
		return true, nil
	}

	logger.WithField("id", id).Errorf("FetchArtistMeta: Artist was not updated but no error was thrown")
	return false, nil
}
