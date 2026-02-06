package service

import (
	"context"
	"fmt"
	infra "github.com/saleh-ghazimoradi/TeleGopher/infra/TXManager"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
	"time"
)

type MessageService interface {
	SendMessage(ctx context.Context, req *dto.CreateMessageRequest, senderId int64) (*dto.MessageResponse, error)
	GetMessage(ctx context.Context, messageId, userId int64) (*dto.MessageResponse, error)
	GetPrivateMessages(ctx context.Context, privateId, userId int64, page, limit int) (*dto.MessagesListResponse, error)
	GetUndeliveredMessages(ctx context.Context, privateId, userId int64) ([]dto.MessageResponse, error)
	MarkMessageAsRead(ctx context.Context, messageId, userId int64) error
	MarkMessageAsDelivered(ctx context.Context, messageId, userId int64) error
}

type messageService struct {
	messageRepository repository.MessageRepository
	privateRepository repository.PrivateRepository
	tx                infra.TxManager
}

func (m *messageService) SendMessage(ctx context.Context, req *dto.CreateMessageRequest, senderId int64) (*dto.MessageResponse, error) {
	private, err := m.privateRepository.GetPrivateById(ctx, req.PrivateId)
	if err != nil {
		return nil, fmt.Errorf("private chat not found: %w", err)
	}

	if private.User1Id != senderId && private.User2Id != senderId {
		return nil, fmt.Errorf("unauthorized to send message in this chat")
	}

	message := &domain.Message{
		FromId:      senderId,
		PrivateId:   req.PrivateId,
		MessageType: domain.MessageType(req.MessageType),
		Content:     req.Content,
		Delivered:   false,
		Read:        false,
		CreatedAt:   time.Now(),
	}

	if err := m.messageRepository.CreateMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return m.toMessageResponse(message), nil
}

func (m *messageService) GetMessage(ctx context.Context, messageId, userId int64) (*dto.MessageResponse, error) {
	message, err := m.messageRepository.GetMessageById(ctx, messageId)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}

	private, err := m.privateRepository.GetPrivateById(ctx, message.PrivateId)
	if err != nil {
		return nil, fmt.Errorf("private chat not found: %w", err)
	}

	if private.User1Id != userId && private.User2Id != userId {
		return nil, fmt.Errorf("unauthorized to view this message")
	}

	return m.toMessageResponse(message), nil
}

func (m *messageService) GetPrivateMessages(ctx context.Context, privateId, userId int64, page, limit int) (*dto.MessagesListResponse, error) {
	private, err := m.privateRepository.GetPrivateById(ctx, privateId)
	if err != nil {
		return nil, fmt.Errorf("private chat not found: %w", err)
	}

	if private.User1Id != userId && private.User2Id != userId {
		return nil, fmt.Errorf("unauthorized to view messages in this chat")
	}

	messages, err := m.messageRepository.GetMessageByPrivateId(ctx, privateId, page, limit+1)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	hasNextPage := false
	if len(messages) > limit {
		hasNextPage = true
		messages = messages[:limit]
	}

	response := &dto.MessagesListResponse{
		Messages:    make([]dto.MessageResponse, len(messages)),
		Page:        page,
		Limit:       limit,
		HasNextPage: hasNextPage,
	}

	for i, msg := range messages {
		response.Messages[i] = *m.toMessageResponse(&msg)
	}

	return response, nil
}

func (m *messageService) GetUndeliveredMessages(ctx context.Context, privateId, userId int64) ([]dto.MessageResponse, error) {
	private, err := m.privateRepository.GetPrivateById(ctx, privateId)
	if err != nil {
		return nil, fmt.Errorf("private chat not found: %w", err)
	}

	if private.User1Id != userId && private.User2Id != userId {
		return nil, fmt.Errorf("unauthorized to view messages in this chat")
	}

	messages, err := m.messageRepository.GetUndeliveredMessagesByPrivateId(ctx, privateId)
	if err != nil {
		return nil, fmt.Errorf("failed to get undelivered messages: %w", err)
	}

	response := make([]dto.MessageResponse, len(messages))
	for i, msg := range messages {
		response[i] = *m.toMessageResponse(&msg)
	}

	return response, nil
}

func (m *messageService) MarkMessageAsRead(ctx context.Context, messageId, userId int64) error {
	message, err := m.messageRepository.GetMessageById(ctx, messageId)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	private, err := m.privateRepository.GetPrivateById(ctx, message.PrivateId)
	if err != nil {
		return fmt.Errorf("private chat not found: %w", err)
	}

	var recipientId int64
	if private.User1Id == message.FromId {
		recipientId = private.User2Id
	} else {
		recipientId = private.User1Id
	}

	if recipientId != userId {
		return fmt.Errorf("unauthorized to mark this message as read")
	}

	return m.messageRepository.MarkMessageRead(ctx, messageId)
}

func (m *messageService) MarkMessageAsDelivered(ctx context.Context, messageId, userId int64) error {
	message, err := m.messageRepository.GetMessageById(ctx, messageId)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	private, err := m.privateRepository.GetPrivateById(ctx, message.PrivateId)
	if err != nil {
		return fmt.Errorf("private chat not found: %w", err)
	}

	var recipientId int64
	if private.User1Id == message.FromId {
		recipientId = private.User2Id
	} else {
		recipientId = private.User1Id
	}

	if recipientId != userId {
		return fmt.Errorf("unauthorized to mark this message as delivered")
	}

	return m.messageRepository.MarkMessageDelivered(ctx, messageId)
}

func (m *messageService) toMessageResponse(message *domain.Message) *dto.MessageResponse {
	return &dto.MessageResponse{
		Id:          message.Id,
		FromId:      message.FromId,
		PrivateId:   message.PrivateId,
		MessageType: string(message.MessageType),
		Content:     message.Content,
		Delivered:   message.Delivered,
		Read:        message.Read,
		CreatedAt:   message.CreatedAt,
	}
}

func NewMessageService(messageRepository repository.MessageRepository, privateRepository repository.PrivateRepository, tx infra.TxManager) MessageService {
	return &messageService{
		messageRepository: messageRepository,
		privateRepository: privateRepository,
		tx:                tx,
	}
}
