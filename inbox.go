package evilmail

import (
	"context"
	"fmt"
	"net/url"
)

// InboxService handles operations related to account inboxes and messages.
//
// This service is for reading messages from persistent email accounts.
// For temporary email messages, use TempEmailService.GetMessage instead.
type InboxService struct {
	client *Client
}

// List retrieves all messages in the inbox for the specified email address.
// Returns a list of message summaries (without full body content).
//
//	messages, err := client.Inbox.List(ctx, "hello@evilmail.pro")
//	for _, msg := range messages {
//	    fmt.Printf("[%s] %s: %s\n", msg.UID, msg.From, msg.Subject)
//	}
func (s *InboxService) List(ctx context.Context, email string) ([]MessageSummary, error) {
	if email == "" {
		return nil, &ValidationError{Field: "email", Message: "email is required"}
	}

	path := fmt.Sprintf("/api/ext/accounts/inbox?email=%s", url.QueryEscape(email))
	return getJSON[[]MessageSummary](s.client, ctx, path)
}

// GetMessage retrieves the full content of a specific message by its UID.
// The email parameter identifies which account's mailbox to read from.
//
//	msg, err := client.Inbox.GetMessage(ctx, "12345", "hello@evilmail.pro")
//	fmt.Println(msg.Subject)
//	fmt.Println(msg.Text)
func (s *InboxService) GetMessage(ctx context.Context, uid, email string) (*Message, error) {
	if uid == "" {
		return nil, &ValidationError{Field: "uid", Message: "message UID is required"}
	}
	if email == "" {
		return nil, &ValidationError{Field: "email", Message: "email is required"}
	}

	path := fmt.Sprintf("/api/ext/accounts/message/%s?email=%s", url.PathEscape(uid), url.QueryEscape(email))
	result, err := getJSON[Message](s.client, ctx, path)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
