package services

import (
	irepository "github.com/ESSantana/streaming-test/internal/repositories/interfaces"
	iservices "github.com/ESSantana/streaming-test/internal/services/interfaces"
	istorage "github.com/ESSantana/streaming-test/internal/storage/interfaces"
)

type serviceManager struct {
	storageManager    istorage.StorageManager
	repositoryManager irepository.RepositoryManager
}

func NewServiceManager(storageManager istorage.StorageManager, repositoryManager irepository.RepositoryManager) iservices.ServiceManager {
	return &serviceManager{
		storageManager:    storageManager,
		repositoryManager: repositoryManager,
	}
}

func (manager *serviceManager) NewVideoService() iservices.VideoService {
	return newVideoService(manager.storageManager, manager.repositoryManager)
}
