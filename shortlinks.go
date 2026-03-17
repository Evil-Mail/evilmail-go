package evilmail

import "context"

// ShortlinksService handles operations related to shortlink creation.
//
// Shortlinks create short-lived URLs for sharing temporary email messages.
// They are stored in Redis and expire when the associated temp email session expires.
type ShortlinksService struct {
	client *Client
}

// Create creates a new shortlink for a temporary email session or message.
// The request must include a session token and a type ("read" or "open").
// Optionally include a message UID to link to a specific message.
//
//	link, err := client.Shortlinks.Create(ctx, &evilmail.ShortlinkRequest{
//	    Token: temp.SessionToken,
//	    UID:   "12345",
//	    Type:  "read",
//	})
//	fmt.Printf("Shortlink: %s\n", link.URL)
func (s *ShortlinksService) Create(ctx context.Context, req *ShortlinkRequest) (*Shortlink, error) {
	if req == nil {
		return nil, &ValidationError{Field: "request", Message: "request body is required"}
	}
	if req.Token == "" {
		return nil, &ValidationError{Field: "token", Message: "session token is required"}
	}
	if req.Type == "" {
		return nil, &ValidationError{Field: "type", Message: "shortlink type is required"}
	}
	if req.Type != "read" && req.Type != "open" {
		return nil, &ValidationError{Field: "type", Message: "shortlink type must be \"read\" or \"open\""}
	}

	result, err := postJSON[Shortlink](s.client, ctx, "/api/ext/shortlink", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
