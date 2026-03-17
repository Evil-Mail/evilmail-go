package evilmail

import "context"

// RandomEmailService handles operations related to random email generation.
//
// This service provides both preview (without creation) and batch creation
// of random email accounts.
type RandomEmailService struct {
	client *Client
}

// Preview generates a random email preview without creating an account.
// This is useful for previewing what a generated email address would look like.
//
//	preview, err := client.RandomEmail.Preview(ctx)
//	fmt.Printf("Generated: %s (password: %s)\n", preview.Email, preview.Password)
func (s *RandomEmailService) Preview(ctx context.Context) (*RandomEmailPreview, error) {
	result, err := getJSON[RandomEmailPreview](s.client, ctx, "/api/random-email")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchCreate creates multiple random email accounts at once.
// The domain field is required; count and password length are optional
// and will use server defaults if set to zero.
//
//	resp, err := client.RandomEmail.BatchCreate(ctx, &evilmail.BatchCreateRandomEmailRequest{
//	    Domain:         "evilmail.pro",
//	    Count:          5,
//	    PasswordLength: 16,
//	})
//	for _, cred := range resp.Emails {
//	    fmt.Printf("%s : %s\n", cred.Email, cred.Password)
//	}
func (s *RandomEmailService) BatchCreate(ctx context.Context, req *BatchCreateRandomEmailRequest) (*BatchCreateRandomEmailResponse, error) {
	if req == nil {
		return nil, &ValidationError{Field: "request", Message: "request body is required"}
	}
	if req.Domain == "" {
		return nil, &ValidationError{Field: "domain", Message: "domain is required"}
	}

	result, err := postJSON[BatchCreateRandomEmailResponse](s.client, ctx, "/api/random-email", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
