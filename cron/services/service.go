package services

import "nft-raffle-cron/database"

type Container struct {
	HelloService *HelloService

	NftRaffleMongoDb *database.NftRaffleMongoDb
}

func NewContainer() *Container {
	nftRaffleMongoDb := database.GetNftRaffleMongoDb()

	return &Container{
		HelloService: GetHelloService(nftRaffleMongoDb),

		NftRaffleMongoDb: nftRaffleMongoDb,
	}
}
