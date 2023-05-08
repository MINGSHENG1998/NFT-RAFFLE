package services

import (
	"context"
	"fmt"
	"nft-raffle-cron/database"
	"nft-raffle-cron/logger"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
)

const USED_REFRESH_TOKEN = "usedRefreshToken"

var (
	usedRefreshTokenService     *UsedRefreshTokenService
	usedRefreshTokenServiceOnce sync.Once
)

type UsedRefreshTokenService struct {
	nftRaffleMongoDb *database.NftRaffleMongoDb
}

func GetUsedRefreshTokenService(nftRaffleMongoDb *database.NftRaffleMongoDb) *UsedRefreshTokenService {
	if usedRefreshTokenService == nil {
		usedRefreshTokenServiceOnce.Do(func() {
			usedRefreshTokenService = &UsedRefreshTokenService{
				nftRaffleMongoDb: nftRaffleMongoDb,
			}
		})
	}
	return usedRefreshTokenService
}

func (s *UsedRefreshTokenService) CountAllUsedRefreshToken() (int64, error) {
	usedRefreshTokenCollection := s.nftRaffleMongoDb.OpenCollection(s.nftRaffleMongoDb.GetClient(), USED_REFRESH_TOKEN)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return usedRefreshTokenCollection.CountDocuments(ctx, bson.M{})
}

func (s *UsedRefreshTokenService) TestCount() {
	loc, err := timeUtil.GetCurrentLocation()
	if err != nil {
		logger.Logger.Panic("unable to load current location")
	}
	scheduler := gocron.NewScheduler(loc)
	scheduler.Every(3).Seconds().Do(func() {
		count, err := s.CountAllUsedRefreshToken()
		if err != nil {
			logger.Logger.Warn(fmt.Sprintf("error occured within TestCount(): %v", err.Error()))
			return
		}
		logger.Logger.Info(fmt.Sprintf("Used refresh token count: %v", count))
	})
	scheduler.StartAsync()
}
