package services

import "nft-raffle-cron/database"

type Container struct {
	HelloService            *HelloService
	UsedRefreshTokenService *UsedRefreshTokenService

	NftRaffleMongoDb *database.NftRaffleMongoDb
}

func NewContainer() *Container {
	nftRaffleMongoDb := database.GetNftRaffleMongoDb()

	return &Container{
		HelloService:            GetHelloService(nftRaffleMongoDb),
		UsedRefreshTokenService: GetUsedRefreshTokenService(nftRaffleMongoDb),

		NftRaffleMongoDb: nftRaffleMongoDb,
	}
}
