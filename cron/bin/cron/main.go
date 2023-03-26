package main

import (
	"fmt"
	"nft-raffle-cron/services"
)

func main() {
	forever := make(chan struct{})

	serviceContainer := services.NewServiceContainer()
	serviceContainer.HelloService.HelloCron()

	fmt.Println("Press ctrl+C to exit")
	<-forever
}
