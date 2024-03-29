package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/Denis-Carlos-Farias/upload-S3/cmd/service"
	"github.com/joho/godotenv"
)

var (
	wg sync.WaitGroup
)

func init() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load("../config/.env")
	if err != nil {
		log.Fatal("An error has occurred to get the file .env", err)
	}

}

func main() {

	sc := service.NewProductRepository(&wg)

	dir, err := os.Open(os.Getenv("TARGET_FOLDER"))
	if err != nil {
		panic(err)
	}

	sizechan, err := strconv.Atoi(os.Getenv("CHANEL_SIZE"))
	if err != nil {
		panic(err)
	}

	trafficLight := make(chan struct{}, sizechan)

	for {
		files, err := dir.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error: %v\n", err)
			continue
		}

		wg.Add(1)

		trafficLight <- struct{}{}
		go sc.Upload(files[0], trafficLight)
	}

	wg.Wait()
}
