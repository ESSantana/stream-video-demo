package interfaces

type ServiceManager interface { 
	NewVideoService() VideoService
}