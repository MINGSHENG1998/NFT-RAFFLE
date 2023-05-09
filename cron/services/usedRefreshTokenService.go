package services

import (
	"context"
	"fmt"
	"nft-raffle-cron/database"
	"nft-raffle-cron/logger"
	"nft-raffle-cron/utils"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
)

const USED_REFRESH_TOKEN = "usedRefreshToken"

var (
	usedRefreshTokenService     *UsedRefreshTokenService
	usedRefreshTokenServiceOnce sync.Once

	timeUtil = utils.GetTimeUtil()
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

// func (s *UsedRefreshTokenService) CountAllUsedRefreshToken() (int64, error) {
// 	usedRefreshTokenCollection := s.nftRaffleMongoDb.OpenCollection(s.nftRaffleMongoDb.GetClient(), USED_REFRESH_TOKEN)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	return usedRefreshTokenCollection.CountDocuments(ctx, bson.M{})
// }

// func (s *UsedRefreshTokenService) TestCount() {
// 	loc, err := timeUtil.GetCurrentLocation()
// 	if err != nil {
// 		logger.Logger.Panic("unable to load current location")
// 	}
// 	scheduler := gocron.NewScheduler(loc)
// 	scheduler.Every(3).Seconds().Do(func() {
// 		count, err := s.CountAllUsedRefreshToken()
// 		if err != nil {
// 			logger.Logger.Warn(fmt.Sprintf("error occured within TestCount(): %v", err.Error()))
// 			return
// 		}
// 		logger.Logger.Info(fmt.Sprintf("Used refresh token count: %v", count))
// 	})
// 	scheduler.StartAsync()
// }

func (s *UsedRefreshTokenService) StartRemovingUsedRefreshTokenCronAsync() {
	loc, err := timeUtil.GetCurrentLocation()
	if err != nil {
		logger.Logger.Panic("unable to load current location")
	}
	scheduler := gocron.NewScheduler(loc)
	scheduler.Every(1).Day().At("00:05").Do(s.RemoveExpiredUsedRefreshToken)
	scheduler.StartAsync()
}

func (s *UsedRefreshTokenService) RemoveExpiredUsedRefreshToken() {
	usedRefreshTokenCollection := s.nftRaffleMongoDb.OpenCollection(s.nftRaffleMongoDb.GetClient(), USED_REFRESH_TOKEN)

	nowUnix, err := timeUtil.Now()

	if err != nil {
		logger.Logger.Panic("unable to get current now")
		return
	}

	logger.Logger.Info(fmt.Sprintf("now unix: %v", nowUnix))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "expired_at_unix", Value: bson.D{
		{Key: "$lte", Value: nowUnix},
	}}}

	_, err = usedRefreshTokenCollection.DeleteMany(ctx, filter)

	if err != nil {
		logger.Logger.Warn(fmt.Sprintf("error occured when deleting expired used refresh token: %v", err.Error()))
		return
	}

	logger.Logger.Info(fmt.Sprintf("RemoveExpiredUsedRefreshToken function ran for time unix: %v", nowUnix))
}
