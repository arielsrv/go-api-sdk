package services

type IMessageService interface {
	GetMessage() string
}

type MessageService struct {
}

func (r MessageService) GetMessage() string {
	return "Hello World!"
}

func NewMessageService() *MessageService {
	return &MessageService{}
}
