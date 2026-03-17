package evilmail

import (
	"context"
	"fmt"
	"net/url"
)

// TempEmailService handles operations related to temporary email addresses.
//
// Temporary emails are short-lived addresses useful for one-time signups,
// testing, and other scenarios where a disposable email is needed.
type TempEmailService struct {
	client *Client
}

// Create creates a new temporary email address.
//
// The request parameter is optional -- pass nil to use server defaults
// for domain and TTL. Returns the newly created temporary email with
// its session token, which is required for subsequent operations.
//
//	temp, err := client.TempEmail.Create(ctx, &evilmail.CreateTempEmailRequest{
//	    Domain:     "evilmail.pro",
//	    TTLMinutes: 60,
//	})
func (s *TempEmailService) Create(ctx context.Context, req *CreateTempEmailRequest) (*TempEmail, error) {
	if req == nil {
		req = &CreateTempEmailRequest{}
	}

	result, err := postJSON[TempEmail](s.client, ctx, "/api/ext/temp-email", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSession retrieves the session information for a temporary email address.
// Returns the email address, domain, TTL, and expiration time.
// The token is the session token returned by Create.
//
//	session, err := client.TempEmail.GetSession(ctx, temp.SessionToken)
//	fmt.Printf("Email: %s, expires at: %s\n", session.Email, session.ExpiresAt)
func (s *TempEmailService) GetSession(ctx context.Context, token string) (*TempSession, error) {
	if token == "" {
		return nil, &ValidationError{Field: "token", Message: "session token is required"}
	}

	path := fmt.Sprintf("/api/ext/temp-email?token=%s", url.QueryEscape(token))
	result, err := getJSON[TempSession](s.client, ctx, path)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMessage retrieves a specific message from a temporary email inbox by UID.
// The token is the session token returned by Create.
//
//	msg, err := client.TempEmail.GetMessage(ctx, temp.SessionToken, "12345")
//	fmt.Printf("From: %s\nSubject: %s\n", msg.From, msg.Subject)
func (s *TempEmailService) GetMessage(ctx context.Context, token, uid string) (*Message, error) {
	if token == "" {
		return nil, &ValidationError{Field: "token", Message: "session token is required"}
	}
	if uid == "" {
		return nil, &ValidationError{Field: "uid", Message: "message UID is required"}
	}

	path := fmt.Sprintf("/api/ext/message/%s?token=%s", url.PathEscape(uid), url.QueryEscape(token))
	result, err := getJSON[Message](s.client, ctx, path)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a temporary email session and its associated email account.
// This cannot be undone. The token is the session token returned by Create.
//
//	err := client.TempEmail.Delete(ctx, temp.SessionToken)
func (s *TempEmailService) Delete(ctx context.Context, token string) error {
	if token == "" {
		return &ValidationError{Field: "token", Message: "session token is required"}
	}

	path := fmt.Sprintf("/api/ext/temp-email?token=%s", url.QueryEscape(token))
	return s.client.del(ctx, path, nil, nil)
}
