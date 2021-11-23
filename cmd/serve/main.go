package main

import (
	"github.com/spf13/viper"
	"log"
	"micropairs/configs"
	"micropairs/internal/server"
	"net/http"
	"sync"
)

func main() {
	viper.Set(configs.CFGPATH,"../../configs/" )
	if err := configs.InitConfig(); err != nil{
		log.Fatal("Can't read config. ", err)
	}

	s := server.NewCryptoServer(http.DefaultClient)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := s.Run()
		if err != nil{
			log.Fatal("Server error", err)
		}
	}()
	wg.Wait()
}
