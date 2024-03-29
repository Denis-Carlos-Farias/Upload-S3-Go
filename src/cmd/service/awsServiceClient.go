package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/Denis-Carlos-Farias/upload-S3/cmd/domain/interfaces"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ServiceClint struct {
	client *s3.S3
	wg     *sync.WaitGroup
}

func NewProductRepository(w *sync.WaitGroup) interfaces.IServiceClient {

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY"),
			os.Getenv("AWS_SECRET_KEY"), ""),
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		panic(err)
	}

	return &ServiceClint{
		client: s3.New(sess),
		wg:     w,
	}
}

func (sc *ServiceClint) Upload(fileInfo os.FileInfo, trafficLight <-chan struct{}) {
	defer sc.wg.Done()

	fmt.Printf("upload started: %s", fileInfo.Name())
	filePath := fmt.Sprintf("%s/%s", os.Getenv("TARGET_FOLDER"), fileInfo.Name())

	file, err := os.Open(filePath)
	if err != nil {
		<-trafficLight
		fmt.Printf("Error: %v", err)
		return
	}

	defer file.Close()

	var fileSize int64 = fileInfo.Size()

	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	_, err = sc.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("teste"),
		Key:    aws.String(fileInfo.Name()),
		Body:   bytes.NewReader(fileBuffer),
	})
	if err != nil {
		<-trafficLight
		fmt.Printf("Error upload: %v", err)
		return
	}

	fmt.Printf("upload finished: %s", fileInfo.Name())
	<-trafficLight
}
