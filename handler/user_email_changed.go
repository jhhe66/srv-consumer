package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rafaeljesus/srv-consumer"
	"github.com/rafaeljesus/srv-consumer/platform/message"
)

type (
	// UserEmailChanged is the message handler.
	UserEmailChanged struct {
		store srv.UserStore
	}
)

// NewUserEmailChanged returns new UserEmailChanged struct.
func NewUserEmailChanged(s srv.UserStore) *UserEmailChanged {
	return &UserEmailChanged{s}
}

// Handle is the user email changed message handler.
func (u *UserEmailChanged) Handle(ctx context.Context, m *message.Message) error {
	user := new(srv.User)
	if err := json.Unmarshal(m.Body, user); err != nil {
		log.Printf("failed to unmarshal message body: %v", err)
		if err := m.Ack(false); err != nil {
			log.Printf("failed to ack message: %v", err)
		}
		return err
	}

	err := u.store.Save(user)

	switch err {
	case nil:
		log.Print("user status successfully changed")
		if err := m.Ack(false); err != nil {
			return fmt.Errorf("failed to ack message: %v", err)
		}
		return nil
	case srv.ErrNotFound:
		log.Print("user not found")
		if err := m.Ack(false); err != nil {
			log.Printf("failed to ack message: %v", err)
		}
		return err
	default:
		log.Printf("failed to save user to store: %v", err)
		if err := m.Nack(false, true); err != nil {
			log.Printf("failed to reject message: %v", err)
		}
		return err
	}
}
