package services

import (
	"nft-raffle-cron/logger"
	"nft-raffle-cron/utils"
	"sync"

	"github.com/go-co-op/gocron"
)

var (
	helloService     *HelloService
	helloServiceOnce sync.Once

	timeUtil = utils.GetTimeUtil()
)

type HelloService struct{}

func GetHelloService() *HelloService {
	if helloService == nil {
		helloServiceOnce.Do(func() {
			helloService = &HelloService{}
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
