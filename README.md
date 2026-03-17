<p align="center">
  <a href="https://evilmail.pro">
    <img src="https://avatars.githubusercontent.com/u/267867069?v=4" alt="EvilMail Logo" width="120" height="120" style="border-radius: 20px;">
  </a>
</p>

<h1 align="center">EvilMail Go SDK</h1>

<p align="center">
  <strong>Official Go client library for the <a href="https://evilmail.pro">EvilMail</a> disposable email API</strong>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/Evil-Mail/evilmail-go"><img src="https://pkg.go.dev/badge/github.com/Evil-Mail/evilmail-go.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/Evil-Mail/evilmail-go"><img src="https://goreportcard.com/badge/github.com/Evil-Mail/evilmail-go" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License: MIT"></a>
  <a href="https://github.com/Evil-Mail/evilmail-go"><img src="https://img.shields.io/github/go-mod/go-version/Evil-Mail/evilmail-go?style=flat-square" alt="Go Version"></a>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#api-reference">API Reference</a> •
  <a href="#error-handling">Error Handling</a> •
  <a href="https://evilmail.pro/docs">Documentation</a>
</p>

---

The **EvilMail Go SDK** provides a clean, idiomatic Go interface for integrating temporary email, disposable email addresses, email verification code extraction, inbox management, and custom domain email services into your Go applications. Zero external dependencies — built entirely on the Go standard library.

## Features

- **Zero Dependencies** — Built on `net/http` and the Go standard library only
- **Temporary Email** — Create anonymous disposable email addresses with configurable TTL
- **Email Verification Codes** — Auto-extract OTP codes from Google, Facebook, Instagram, TikTok, Discord, Twitter, LinkedIn, iCloud
- **Account Management** — Full CRUD for persistent email accounts on custom domains
- **Inbox Access** — Read emails, list messages, fetch full HTML & plain text content
- **Random Email Generator** — Batch create random email accounts with auto-generated passwords
- **Domain Management** — List free, premium, and custom email domains
- **Shortlink Creation** — Generate short URLs for temporary email sessions
- **Context Support** — Full `context.Context` propagation for cancellation and deadlines
- **Functional Options** — Clean, extensible client configuration pattern
- **Typed Errors** — Structured error types with `errors.As` / `errors.Is` support
- **Thread Safe** — Safe for concurrent use from multiple goroutines

## Requirements

- Go 1.21 or later
- An EvilMail API key — [Get yours free](https://evilmail.pro)

## Installation

```bash
go get github.com/Evil-Mail/evilmail-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	evilmail "github.com/Evil-Mail/evilmail-go"
)

func main() {
	client := evilmail.New("your-api-key")
	ctx := context.Background()

	// Create a temporary disposable email address
	temp, err := client.TempEmail.Create(ctx, &evilmail.CreateTempEmailRequest{
		Domain:     "evilmail.pro",
		TTLMinutes: 60,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Temporary email: %s\n", temp.Email)
	fmt.Printf("Session token: %s\n", temp.SessionToken)

	// Check session status
	session, err := client.TempEmail.GetSession(ctx, temp.SessionToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Expires at: %s\n", session.ExpiresAt)

	// Read a specific message from temp inbox
	msg, err := client.TempEmail.GetMessage(ctx, temp.SessionToken, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Subject: %s\n", msg.Subject)

	// Extract a Google verification code
	code, err := client.Verification.GetCode(ctx, evilmail.ServiceGoogle, "user@yourdomain.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Verification code: %s\n", code.Code)

	// List all accounts
	accounts, err := client.Accounts.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, acct := range accounts {
		fmt.Printf("Account: %s\n", acct.Email)
	}

	// Batch create random email accounts
	batch, err := client.RandomEmail.BatchCreate(ctx, &evilmail.BatchCreateRandomEmailRequest{
		Domain:         "yourdomain.com",
		Count:          5,
		PasswordLength: 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created %d random accounts\n", batch.Count)
}
```

## Client Configuration

The client uses the functional options pattern for clean, extensible configuration:

```go
import "time"

// Default configuration
client := evilmail.New("your-api-key")

// Custom configuration
client := evilmail.New("your-api-key",
	evilmail.WithBaseURL("https://your-instance.example.com"),
	evilmail.WithTimeout(60 * time.Second),
	evilmail.WithHTTPClient(myCustomHTTPClient),
)
```

| Option | Description |
|--------|-------------|
| `WithBaseURL(url)` | Override the default API base URL |
| `WithTimeout(d)` | Set HTTP request timeout (default: 30s) |
| `WithHTTPClient(c)` | Provide a custom `*http.Client` |

---

## API Reference

All methods accept a `context.Context` as their first argument for cancellation, deadlines, and tracing propagation.

### Temporary Email

Create anonymous, disposable email addresses with automatic expiration. Ideal for sign-up verification, testing, and privacy protection.

#### Create Temporary Email

```go
temp, err := client.TempEmail.Create(ctx, &evilmail.CreateTempEmailRequest{
	Domain:     "evilmail.pro", // optional
	TTLMinutes: 30,             // optional, minutes until auto-expiry
})
```

Returns: `*TempEmail` — `Email`, `Domain`, `SessionToken`, `TTLMinutes`, `ExpiresAt`

#### Get Session Status

```go
session, err := client.TempEmail.GetSession(ctx, token)
```

Check if a temporary email session is still active and retrieve session details.

#### Get Message

```go
msg, err := client.TempEmail.GetMessage(ctx, token, uid)
```

Read a specific message from a temporary email inbox.

#### Delete Session

```go
err := client.TempEmail.Delete(ctx, token)
```

Permanently delete a temporary email session and all associated data.

---

### Accounts

Manage persistent email accounts on custom domains.

```go
// List all accounts
accounts, err := client.Accounts.List(ctx)

// Create a new account
resp, err := client.Accounts.Create(ctx, &evilmail.CreateAccountRequest{
	Email:    "user@yourdomain.com",
	Password: "secure-password",
})

// Delete accounts
deleted, err := client.Accounts.Delete(ctx, &evilmail.DeleteAccountsRequest{
	Emails: []string{"old@yourdomain.com"},
})

// Change password
err = client.Accounts.ChangePassword(ctx, &evilmail.ChangePasswordRequest{
	Email:       "user@yourdomain.com",
	NewPassword: "new-secure-password",
})
```

---

### Inbox & Messages

Read emails from persistent account inboxes with full message content.

```go
// List inbox messages
messages, err := client.Inbox.List(ctx, "user@yourdomain.com")
for _, msg := range messages {
	fmt.Printf("[%d] %s: %s\n", msg.UID, msg.From, msg.Subject)
}

// Read full message content (HTML + plain text)
msg, err := client.Inbox.GetMessage(ctx, 12345, "user@yourdomain.com")
fmt.Printf("HTML: %s\nText: %s\n", msg.HTML, msg.Text)
```

---

### Verification Codes

Automatically extract OTP verification codes from emails sent by popular services.

```go
code, err := client.Verification.GetCode(ctx, evilmail.ServiceGoogle, "user@yourdomain.com")
fmt.Printf("Code: %s (from: %s)\n", code.Code, code.From)
```

**Supported services:**

| Constant | Service |
|----------|---------|
| `ServiceFacebook` | Facebook |
| `ServiceTwitter` | Twitter / X |
| `ServiceGoogle` | Google |
| `ServiceICloud` | iCloud |
| `ServiceInstagram` | Instagram |
| `ServiceTikTok` | TikTok |
| `ServiceDiscord` | Discord |
| `ServiceLinkedIn` | LinkedIn |

---

### Random Email

Generate random email accounts with secure auto-generated credentials.

```go
// Preview a random email (without creating it)
preview, err := client.RandomEmail.Preview(ctx)
fmt.Printf("%s : %s\n", preview.Email, preview.Password)

// Batch create random email accounts
batch, err := client.RandomEmail.BatchCreate(ctx, &evilmail.BatchCreateRandomEmailRequest{
	Domain:         "yourdomain.com",
	Count:          10,
	PasswordLength: 20,
})
for _, cred := range batch.Emails {
	fmt.Printf("%s : %s\n", cred.Email, cred.Password)
}
```

---

### Domains

List available email domains by tier.

```go
// List customer domains (authenticated)
domains, err := client.Domains.List(ctx)
fmt.Println("Free:", domains.Free)
fmt.Println("Premium:", domains.Premium)
fmt.Println("Custom:", domains.Customer)

// List public domains (unauthenticated)
public, err := client.PublicDomains.List(ctx)
fmt.Println("Domains:", public.Domains)
```

---

### Shortlinks

Generate short URLs for temporary email sessions.

```go
link, err := client.Shortlinks.Create(ctx, &evilmail.ShortlinkRequest{
	Token: sessionToken,
	Type:  "session",
})
fmt.Printf("Short URL: %s\n", link.ShortURL)
```

---

## Error Handling

The SDK provides typed errors compatible with Go's `errors.As` and `errors.Is` patterns:

```go
import "errors"

msg, err := client.Inbox.GetMessage(ctx, uid, email)
if err != nil {
	// Check for specific API errors
	if evilmail.IsNotFoundError(err) {
		fmt.Println("Message not found")
		return
	}
	if evilmail.IsAuthError(err) {
		fmt.Println("Invalid API key — check your credentials")
		return
	}
	if evilmail.IsRateLimitError(err) {
		fmt.Println("Rate limited — slow down and retry")
		return
	}

	// Generic API error with status code
	var apiErr *evilmail.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
		return
	}

	// Client-side validation error
	var valErr *evilmail.ValidationError
	if errors.As(err, &valErr) {
		fmt.Printf("Invalid parameter %s: %s\n", valErr.Field, valErr.Message)
		return
	}

	// Network or other error
	log.Fatal(err)
}
```

### Error Types

| Type | Description |
|------|-------------|
| `*APIError` | Non-2xx API response with `StatusCode`, `Status`, `Message`, `Body` |
| `*AuthError` | 401/403 — wraps `*APIError` for authentication failures |
| `*ValidationError` | Client-side validation with `Field` and `Message` |

### Helper Functions

| Function | Description |
|----------|-------------|
| `IsNotFoundError(err)` | Returns true for 404 responses |
| `IsAuthError(err)` | Returns true for 401/403 responses |
| `IsRateLimitError(err)` | Returns true for 429 responses |

---

## Context Support

Every method accepts `context.Context` for cancellation, deadlines, and distributed tracing:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
messages, err := client.Inbox.List(ctx, "user@yourdomain.com")

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
	<-sigChan
	cancel()
}()
temp, err := client.TempEmail.Create(ctx, nil)
```

---

## Use Cases

- **Automated Testing** — Generate disposable email addresses for end-to-end test suites
- **Web Scraping & Crawling** — Create temp emails for sign-up verification in automation pipelines
- **Email Verification Bots** — Automatically extract OTP codes from Google, Facebook, Instagram, and more
- **Microservices** — Lightweight email client for Go microservices and serverless functions
- **CLI Tools** — Build email automation scripts and command-line utilities
- **DevOps & CI/CD** — Integrate email testing into CI pipelines (GitHub Actions, GitLab CI)
- **SaaS Backend** — Automate email provisioning and verification in your Go backend
- **Privacy & Anonymity** — Use anonymous temporary email for privacy-sensitive workflows
- **Kubernetes & Cloud Native** — Thread-safe, context-aware client for cloud-native applications

---

## Related SDKs

| Language | Package | Repository |
|----------|---------|------------|
| **Node.js** | `evilmail` | [Evil-Mail/evilmail-node](https://github.com/Evil-Mail/evilmail-node) |
| **PHP** | `evilmail/evilmail-php` | [Evil-Mail/evilmail-php](https://github.com/Evil-Mail/evilmail-php) |
| **Python** | `evilmail` | [Evil-Mail/evilmail-python](https://github.com/Evil-Mail/evilmail-python) |
| **Go** | `evilmail-go` | [Evil-Mail/evilmail-go](https://github.com/Evil-Mail/evilmail-go) |

## Links

- [EvilMail Website](https://evilmail.pro) — Temporary & custom domain email platform
- [API Documentation](https://evilmail.pro/docs) — Full REST API reference
- [Chrome Extension](https://github.com/Evil-Mail/evilmail-chrome) — Disposable temp email in your browser
- [Firefox Add-on](https://github.com/Evil-Mail/evilmail-firefox) — Temp email for Firefox
- [Mobile App](https://github.com/Evil-Mail/evilmail-mobile) — Privacy-first email on Android

## License

[MIT](LICENSE)

## Support

- Issues: [github.com/Evil-Mail/evilmail-go/issues](https://github.com/Evil-Mail/evilmail-go/issues)
- Email: support@evilmail.pro
- Website: [evilmail.pro](https://evilmail.pro)
