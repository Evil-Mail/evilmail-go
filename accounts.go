package evilmail

import "context"

// AccountsService handles operations related to email accounts.
//
// Accounts are persistent email addresses that can receive and store
// messages until explicitly deleted.
type AccountsService struct {
	client *Client
}

// List retrieves all email accounts associated with the API key.
//
//	accounts, err := client.Accounts.List(ctx)
//	for _, acct := range accounts {
//	    fmt.Printf("%s (%s)\n", acct.Email, acct.Domain)
//	}
func (s *AccountsService) List(ctx context.Context) ([]Account, error) {
	return getJSON[[]Account](s.client, ctx, "/api/accounts")
}

// Create creates a new email account with the given email and password.
//
//	resp, err := client.Accounts.Create(ctx, &evilmail.CreateAccountRequest{
//	    Email:    "hello@evilmail.pro",
//	    Password: "secure-password-123",
//	})
func (s *AccountsService) Create(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	if req == nil {
		return nil, &ValidationError{Field: "request", Message: "request body is required"}
	}
	if req.Email == "" {
		return nil, &ValidationError{Field: "email", Message: "email is required"}
	}
	if req.Password == "" {
		return nil, &ValidationError{Field: "password", Message: "password is required"}
	}

	result, err := postJSON[CreateAccountResponse](s.client, ctx, "/api/accounts", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes one or more email accounts by their email addresses.
// Returns the number of accounts that were successfully deleted.
//
//	resp, err := client.Accounts.Delete(ctx, &evilmail.DeleteAccountsRequest{
//	    Emails: []string{"old@evilmail.pro", "unused@evilmail.pro"},
//	})
//	fmt.Printf("Deleted %d accounts\n", resp.DeletedCount)
func (s *AccountsService) Delete(ctx context.Context, req *DeleteAccountsRequest) (*DeleteAccountsResponse, error) {
	if req == nil {
		return nil, &ValidationError{Field: "request", Message: "request body is required"}
	}
	if len(req.Emails) == 0 {
		return nil, &ValidationError{Field: "emails", Message: "at least one email is required"}
	}

	result, err := delJSON[DeleteAccountsResponse](s.client, ctx, "/api/accounts", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ChangePassword updates the password for the specified email account.
//
//	err := client.Accounts.ChangePassword(ctx, &evilmail.ChangePasswordRequest{
//	    Email:       "hello@evilmail.pro",
//	    NewPassword: "new-secure-password-456",
//	})
func (s *AccountsService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	if req == nil {
		return &ValidationError{Field: "request", Message: "request body is required"}
	}
	if req.Email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}
	if req.NewPassword == "" {
		return &ValidationError{Field: "newPassword", Message: "new password is required"}
	}

	return s.client.put(ctx, "/api/accounts/password", req, nil)
}
