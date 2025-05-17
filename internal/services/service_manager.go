package services

import (
	iservices "github.com/ESSantana/streaming-test/internal/services/interfaces"
	istorage "github.com/ESSantana/streaming-test/internal/storage/interfaces"
)

type serviceManager struct {
	storageManage istorage.StorageManager
}

func NewServiceManager(storageManage istorage.StorageManager) iservices.ServiceManager {
	return &serviceManager{
		storageManage: storageManage,
	}
}

func (svc *serviceManager) NewVideoService() iservices.VideoService {
	return newVideoService(svc.storageManage)
}
