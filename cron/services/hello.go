package services

import (
	"nft-raffle-cron/database"
	"nft-raffle-cron/logger"
	"sync"

	"github.com/go-co-op/gocron"
)

var (
	helloService     *HelloService
	helloServiceOnce sync.Once
)

type HelloService struct {
	nftRaffleMongoDb *database.NftRaffleMongoDb
}

func GetHelloService(nftRaffleMongoDb *database.NftRaffleMongoDb) *HelloService {
	if helloService == nil {
		helloServiceOnce.Do(func() {
			helloService = &HelloService{
				nftRaffleMongoDb: nftRaffleMongoDb,
			}
		})
	}
	return helloService
}

func (s *HelloService) HelloCron() {
	loc, err := timeUtil.GetCurrentLocation()
	if err != nil {
		logger.Logger.Panic("unable to load current location")
	}
	scheduler := gocron.NewScheduler(loc)
	scheduler.Every(3).Seconds().Do(func() {
		logger.Logger.Debug("Hello world")
	})
	scheduler.StartAsync()
}
