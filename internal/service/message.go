package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, input *dto.MessageRequest, senderId uint) (*dto.MessageResponse, error)
	GetMessage(ctx context.Context, messageId, userId uint) (*dto.MessageResponse, error)
	GetPrivateMessages(ctx context.Context, privateId, userId uint, page, limit int) (*dto.MessageListResponse, error)
	GetUndeliveredMessages(ctx context.Context, privateId, userId uint) ([]dto.MessageResponse, error)
	MarkMessageAsRead(ctx context.Context, messageId, userId uint) error
	MarkMessageAsDelivered(ctx context.Context, messageId, userId uint) error
}

type messageService struct {
	messageRepository repository.MessageRepository
	privateRepository repository.PrivateRepository
}

func (m *messageService) SendMessage(ctx context.Context, input *dto.MessageRequest, senderId uint) (*dto.MessageResponse, error) {
	private, err := m.privateRepository.GetPrivateById(ctx, input.PrivateId)
	if err != nil {
		return nil, fmt.Errorf("private chat not found: %w", err)
	}

	if private.User1Id != senderId && private.User2Id != senderId {
		return nil, fmt.Errorf("unauthorized to send message in this chat")
	}

	message := m.toMessageDomain(input, senderId)

	if err := m.messageRepository.CreateMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return m.toMessageDTO(message), nil
}

func (m *messageService) GetMessage(ctx context.Context, messageId, userId uint) (*dto.MessageResponse, error) {
	message, err := m.messageRepository.GetMessageById(ctx, messageId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	if message.PrivateId == nil {
		return nil, fmt.Errorf("message not associated with a private chat")
	}

	private, err := m.privateRepository.GetPrivateById(ctx, *message.PrivateId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("private chat not found")
		}
		return nil, fmt.Errorf("failed to get private chat: %w", err)
	}

	if private.User1Id != userId && private.User2Id != userId {
		return nil, fmt.Errorf("unauthorized to view this message")
	}

	return m.toMessageDTO(message), nil
}

func (m *messageService) GetPrivateMessages(ctx context.Context, privateId, userId uint, page, limit int) (*dto.MessageListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	private, err := m.privateRepository.GetPrivateById(ctx, privateId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("private chat not found")
		}
		return nil, fmt.Errorf("failed to get private chat: %w", err)
	}

	if private.User1Id != userId && private.User2Id != userId {
		return nil, fmt.Errorf("unauthorized to view messages in this chat")
	}

	offset := (page - 1) * limit

	messages, err := m.messageRepository.GetMessageByPrivateId(ctx, privateId, offset, limit+1)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	hasNextPage := false
	if len(messages) > limit {
		hasNextPage = true
		messages = messages[:limit]
	}

	response := &dto.MessageListResponse{
		Messages:    make([]dto.MessageResponse, len(messages)),
		Page:        page,
		Limit:       limit,
		HasNextPage: hasNextPage,
	}

	for i, msg := range messages {
		response.Messages[i] = *m.toMessageDTO(&msg)
	}

	return response, nil
}

func (m *messageService) GetUndeliveredMessages(ctx context.Context, privateId, userId uint) ([]dto.MessageResponse, error) {
	private, err := m.privateRepository.GetPrivateById(ctx, privateId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("private chat not found")
		}
		return nil, fmt.Errorf("failed to get private chat: %w", err)
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
		response[i] = *m.toMessageDTO(&msg)
	}

	return response, nil
}

func (m *messageService) MarkMessageAsRead(ctx context.Context, messageId, userId uint) error {
	message, err := m.messageRepository.GetMessageById(ctx, messageId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("message not found")
		}
		return fmt.Errorf("failed to get message: %w", err)
	}

	if message.PrivateId == nil {
		return fmt.Errorf("message not associated with a private chat")
	}

	private, err := m.privateRepository.GetPrivateById(ctx, *message.PrivateId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("private chat not found")
		}
		return fmt.Errorf("failed to get private chat: %w", err)
	}

	// Determine the recipient (the user who should mark as read)
	var recipientId uint
	if private.User1Id == message.FromId {
		recipientId = private.User2Id
	} else {
		recipientId = private.User1Id
	}

	// Only the recipient can mark the message as read
	if recipientId != userId {
		return fmt.Errorf("only the recipient can mark this message as read")
	}

	if err := m.messageRepository.MarkMessageAsRead(ctx, messageId); err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

func (m *messageService) MarkMessageAsDelivered(ctx context.Context, messageId, userId uint) error {
	message, err := m.messageRepository.GetMessageById(ctx, messageId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("message not found")
		}
		return fmt.Errorf("failed to get message: %w", err)
	}

	if message.PrivateId == nil {
		return fmt.Errorf("message not associated with a private chat")
	}

	private, err := m.privateRepository.GetPrivateById(ctx, *message.PrivateId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("private chat not found")
		}
		return fmt.Errorf("failed to get private chat: %w", err)
	}

	// Determine the recipient
	var recipientId uint
	if private.User1Id == message.FromId {
		recipientId = private.User2Id
	} else {
		recipientId = private.User1Id
	}

	// Only the recipient can mark the message as delivered
	if recipientId != userId {
		return fmt.Errorf("only the recipient can mark this message as delivered")
	}

	if err := m.messageRepository.MarkMessageAsDelivered(ctx, messageId); err != nil {
		return fmt.Errorf("failed to mark message as delivered: %w", err)
	}

	return nil
}

func (m *messageService) toMessageDomain(input *dto.MessageRequest, senderId uint) *domain.Message {
	return &domain.Message{
		FromId:      senderId,
		PrivateId:   &input.PrivateId,
		MessageType: domain.MessageType(input.MessageType),
		Content:     input.Content,
		Delivered:   false,
		Read:        false,
	}
}

func (m *messageService) toMessageDTO(message *domain.Message) *dto.MessageResponse {
	return &dto.MessageResponse{
		Id:          message.Id,
		FromId:      message.FromId,
		MessageType: string(message.MessageType),
		Content:     message.Content,
		Delivered:   message.Delivered,
		Read:        message.Read,
		CreatedAt:   message.CreatedAt,
	}
}

func NewMessageService(messageRepository repository.MessageRepository, privateRepository repository.PrivateRepository) MessageService {
	return &messageService{
		messageRepository: messageRepository,
		privateRepository: privateRepository,
	}
}
