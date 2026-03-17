// Command basic demonstrates the core features of the evilmail-go SDK.
//
// Usage:
//
//	export EVILMAIL_API_KEY="your-api-key"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Evil-Mail/evilmail-go"
)

func main() {
	apiKey := os.Getenv("EVILMAIL_API_KEY")
	if apiKey == "" {
		log.Fatal("EVILMAIL_API_KEY environment variable is required")
	}

	// Create a client with custom timeout.
	client := evilmail.New(apiKey,
		evilmail.WithTimeout(15*time.Second),
	)

	ctx := context.Background()

	// --- Temporary Email ---
	fmt.Println("=== Creating Temporary Email ===")
	temp, err := client.TempEmail.Create(ctx, &evilmail.CreateTempEmailRequest{
		TTLMinutes: 30,
	})
	if err != nil {
		log.Fatalf("Failed to create temp email: %v", err)
	}
	fmt.Printf("Temp email: %s\n", temp.Email)
	fmt.Printf("Session token: %s\n", temp.SessionToken)
	fmt.Printf("Expires at: %s\n", temp.ExpiresAt.Format(time.RFC3339))
	fmt.Println()

	// Check the temporary email session status.
	session, err := client.TempEmail.GetSession(ctx, temp.SessionToken)
	if err != nil {
		log.Fatalf("Failed to get temp session: %v", err)
	}
	fmt.Printf("Temp email session: %s (domain: %s)\n\n", session.Email, session.Domain)

	// --- List Domains ---
	fmt.Println("=== Available Domains ===")
	domains, err := client.Domains.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list domains: %v", err)
	}
	fmt.Printf("Free domains: %v\n", domains.Free)
	fmt.Printf("Premium domains: %v\n", domains.Premium)
	fmt.Printf("Package type: %s\n\n", domains.PackageType)

	// --- Public Domains (no auth required) ---
	fmt.Println("=== Public Domains ===")
	pubDomains, err := client.PublicDomains.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list public domains: %v", err)
	}
	fmt.Printf("Free domains: %v\n", pubDomains.Free)
	fmt.Printf("TTL options: %v\n\n", pubDomains.TTLOptions)

	// --- Random Email Preview ---
	fmt.Println("=== Random Email Preview ===")
	preview, err := client.RandomEmail.Preview(ctx)
	if err != nil {
		log.Fatalf("Failed to preview random email: %v", err)
	}
	fmt.Printf("Username: %s\n", preview.Username)
	fmt.Printf("Email: %s\n", preview.Email)
	fmt.Printf("Password: %s\n\n", preview.Password)

	// --- List Accounts ---
	fmt.Println("=== Listing Accounts ===")
	accounts, err := client.Accounts.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list accounts: %v", err)
	}
	for _, acct := range accounts {
		fmt.Printf("  %s (%s) - created %s\n", acct.Email, acct.Domain, acct.CreatedAt.Format(time.RFC3339))
	}
	if len(accounts) == 0 {
		fmt.Println("  (no accounts)")
	}
	fmt.Println()

	// --- Read Inbox (if accounts exist) ---
	if len(accounts) > 0 {
		email := accounts[0].Email
		fmt.Printf("=== Inbox for %s ===\n", email)
		messages, err := client.Inbox.List(ctx, email)
		if err != nil {
			log.Fatalf("Failed to list inbox: %v", err)
		}
		for _, msg := range messages {
			fmt.Printf("  [%s] %s: %s (seen: %v)\n", msg.UID, msg.From, msg.Subject, msg.Seen)
		}
		if len(messages) == 0 {
			fmt.Println("  (empty inbox)")
		}
		fmt.Println()

		// Read the first message in full.
		if len(messages) > 0 {
			fmt.Printf("=== Reading Message %s ===\n", messages[0].UID)
			msg, err := client.Inbox.GetMessage(ctx, messages[0].UID, email)
			if err != nil {
				log.Fatalf("Failed to read message: %v", err)
			}
			fmt.Printf("From: %s\n", msg.From)
			fmt.Printf("Subject: %s\n", msg.Subject)
			fmt.Printf("Date: %s\n", msg.Date.Format(time.RFC3339))
			fmt.Printf("Text body:\n%s\n", msg.Text)
			fmt.Println()
		}

		// Try to extract a Google verification code.
		fmt.Printf("=== Verification Code (Google) for %s ===\n", email)
		code, err := client.Verification.GetCode(ctx, evilmail.ServiceGoogle, email)
		if err != nil {
			fmt.Printf("  No code found: %v\n", err)
		} else {
			fmt.Printf("  Code: %s (from: %s)\n", code.Code, code.From)
		}
		fmt.Println()
	}

	// --- Cleanup: delete the temp email ---
	fmt.Println("=== Deleting Temporary Email ===")
	if err := client.TempEmail.Delete(ctx, temp.SessionToken); err != nil {
		log.Fatalf("Failed to delete temp email: %v", err)
	}
	fmt.Println("Temp email deleted successfully")
}
