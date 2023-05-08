package main

import (
	"fmt"
	"nft-raffle-cron/services"
)

func main() {
	forever := make(chan struct{})

	container := services.NewContainer()
	container.UsedRefreshTokenService.TestCount()

	fmt.Println("Press ctrl+C to exit")
	<-forever
}
