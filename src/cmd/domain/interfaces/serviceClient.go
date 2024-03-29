package interfaces

import "os"

type IServiceClient interface {
	Upload(fileInfo os.FileInfo, trafficLight <-chan struct{})
}
