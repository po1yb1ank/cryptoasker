package main

import (
	"log"
	"net/http"
	"sync"

	"micropairs/configs"
	"micropairs/internal/server"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		log.Fatal("Can't read config. ", err)
	}

	s := server.NewCryptoServer(http.DefaultClient)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := s.Run()
		if err != nil {
			log.Fatal("Server error", err)
		}
	}()
	wg.Wait()
}
