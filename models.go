package evilmail

import "time"

// apiResponse is the standard envelope returned by all EvilMail API endpoints.
// All successful responses have the shape { "status": "success", "data": ... }.
type apiResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

// TempSession represents the session information for an active temporary email.
// Returned when checking the status of a temporary email session.
type TempSession struct {
	// Email is the full temporary email address.
	Email string `json:"email"`

	// Domain is the domain portion of the email address.
	Domain string `json:"domain"`

	// TTLMinutes is the time-to-live in minutes before the email expires.
	TTLMinutes int `json:"ttlMinutes"`

	// ExpiresAt is the timestamp when the temporary email will expire.
	ExpiresAt time.Time `json:"expiresAt"`
}

// ShortlinkRequest contains the parameters for creating a shortlink.
type ShortlinkRequest struct {
	// Token is the temp email session token. Required.
	Token string `json:"token"`

	// UID is the message UID to link to (optional, depends on type).
	UID string `json:"uid,omitempty"`

	// Type is the shortlink type: "read" or "open". Required.
	Type string `json:"type"`
}

// Shortlink represents a created shortlink.
type Shortlink struct {
	// Code is the short code identifier.
	Code string `json:"code"`

	// URL is the full shortlink URL.
	URL string `json:"url"`
}

// PublicDomains represents the publicly available email domains
// returned by the unauthenticated domains endpoint.
type PublicDomains struct {
	// Free contains the list of free-tier domains.
	Free []string `json:"free"`

	// Premium contains the list of premium domains.
	Premium []string `json:"premium"`

	// TTLOptions contains the available TTL options in minutes.
	TTLOptions []int `json:"ttlOptions"`
}

// CreateTempEmailRequest contains the parameters for creating a temporary email address.
type CreateTempEmailRequest struct {
	// Domain specifies which domain to use for the temporary email.
	// If empty, the server will select a default domain.
	Domain string `json:"domain,omitempty"`

	// TTLMinutes specifies how long the temporary email should remain active,
	// in minutes. If zero, the server will use a default TTL.
	TTLMinutes int `json:"ttlMinutes,omitempty"`
}

// TempEmail represents a newly created temporary email address.
type TempEmail struct {
	// Email is the full temporary email address.
	Email string `json:"email"`

	// Domain is the domain portion of the email address.
	Domain string `json:"domain"`

	// SessionToken is the token used to retrieve the inbox for this
	// temporary email address.
	SessionToken string `json:"sessionToken"`

	// TTLMinutes is the time-to-live in minutes before the email expires.
	TTLMinutes int `json:"ttlMinutes"`

	// ExpiresAt is the timestamp when the temporary email will expire.
	ExpiresAt time.Time `json:"expiresAt"`
}

// MessageSummary represents a brief overview of an email message,
// as returned in inbox listings.
type MessageSummary struct {
	// UID is the unique identifier of the message.
	UID string `json:"uid"`

	// From is the sender's email address.
	From string `json:"from"`

	// Subject is the email subject line.
	Subject string `json:"subject"`

	// Date is the timestamp when the message was received.
	Date time.Time `json:"date"`

	// Seen indicates whether the message has been read.
	Seen bool `json:"seen"`
}

// Account represents an email account.
type Account struct {
	// Email is the full email address.
	Email string `json:"email"`

	// Domain is the domain portion of the email address.
	Domain string `json:"domain"`

	// CreatedAt is the timestamp when the account was created.
	CreatedAt time.Time `json:"createdAt"`
}

// CreateAccountRequest contains the parameters for creating a new email account.
type CreateAccountRequest struct {
	// Email is the desired email address. Required.
	Email string `json:"email"`

	// Password is the password for the new account. Required.
	Password string `json:"password"`
}

// CreateAccountResponse contains the result of creating an email account.
type CreateAccountResponse struct {
	// Email is the email address that was created.
	Email string `json:"email"`
}

// DeleteAccountsRequest contains the parameters for deleting email accounts.
type DeleteAccountsRequest struct {
	// Emails is the list of email addresses to delete.
	Emails []string `json:"emails"`
}

// DeleteAccountsResponse contains the result of a bulk account deletion.
type DeleteAccountsResponse struct {
	// DeletedCount is the number of accounts that were deleted.
	DeletedCount int `json:"deletedCount"`
}

// ChangePasswordRequest contains the parameters for changing an account password.
type ChangePasswordRequest struct {
	// Email is the email address of the account. Required.
	Email string `json:"email"`

	// NewPassword is the new password to set. Required.
	NewPassword string `json:"newPassword"`
}

// Message represents a full email message with body content.
type Message struct {
	// UID is the unique identifier of the message.
	UID string `json:"uid"`

	// From is the sender's email address.
	From string `json:"from"`

	// Subject is the email subject line.
	Subject string `json:"subject"`

	// Text is the plain text body of the message.
	Text string `json:"text"`

	// HTML is the HTML body of the message.
	HTML string `json:"html"`

	// Date is the timestamp when the message was received.
	Date time.Time `json:"date"`

	// Seen indicates whether the message has been read.
	Seen bool `json:"seen"`
}

// VerificationCode represents an extracted verification code from an email.
type VerificationCode struct {
	// Code is the extracted verification code.
	Code string `json:"code"`

	// Service is the service that sent the verification code.
	Service string `json:"service"`

	// Email is the email address that received the code.
	Email string `json:"email"`

	// From is the sender's email address.
	From string `json:"from"`

	// Subject is the subject of the email containing the code.
	Subject string `json:"subject"`

	// Date is the timestamp of the email containing the code.
	Date time.Time `json:"date"`
}

// RandomEmailPreview represents a preview of a randomly generated email account.
type RandomEmailPreview struct {
	// Username is the username portion of the generated email.
	Username string `json:"username"`

	// Email is the full generated email address.
	Email string `json:"email"`

	// Password is the generated password.
	Password string `json:"password"`
}

// BatchCreateRandomEmailRequest contains the parameters for batch creating
// random email accounts.
type BatchCreateRandomEmailRequest struct {
	// Domain is the domain to use for the random emails. Required.
	Domain string `json:"domain"`

	// Count is the number of random emails to create.
	// If zero, the server will use a default count.
	Count int `json:"count,omitempty"`

	// PasswordLength is the length of generated passwords.
	// If zero, the server will use a default length.
	PasswordLength int `json:"passwordLength,omitempty"`
}

// BatchCreateRandomEmailResponse contains the result of a batch random email creation.
type BatchCreateRandomEmailResponse struct {
	// Count is the number of email accounts created.
	Count int `json:"count"`

	// Emails is the list of created email accounts with their credentials.
	Emails []EmailCredentials `json:"emails"`
}

// EmailCredentials represents an email address paired with its password.
type EmailCredentials struct {
	// Email is the email address.
	Email string `json:"email"`

	// Password is the account password.
	Password string `json:"password"`
}

// Domains represents the available email domains grouped by type.
type Domains struct {
	// Free contains the list of free-tier domains.
	Free []string `json:"free"`

	// Premium contains the list of premium domains.
	Premium []string `json:"premium"`

	// Customer contains the list of customer-owned domains.
	Customer []string `json:"customer"`

	// PackageType is the current user's package type.
	PackageType string `json:"packageType"`

	// Authenticated indicates whether the request was made with valid credentials.
	Authenticated bool `json:"authenticated"`
}
