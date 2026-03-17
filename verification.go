package evilmail

import (
	"context"
	"fmt"
	"net/url"
)

// Service name constants for use with VerificationService.GetCode.
// These represent the supported services for verification code extraction.
const (
	ServiceFacebook  = "facebook"
	ServiceTwitter   = "twitter"
	ServiceGoogle    = "google"
	ServiceICloud    = "icloud"
	ServiceInstagram = "instagram"
	ServiceTikTok    = "tiktok"
	ServiceDiscord   = "discord"
	ServiceLinkedIn  = "linkedin"
)

// validServices is the set of supported verification code services.
var validServices = map[string]bool{
	ServiceFacebook:  true,
	ServiceTwitter:   true,
	ServiceGoogle:    true,
	ServiceICloud:    true,
	ServiceInstagram: true,
	ServiceTikTok:    true,
	ServiceDiscord:   true,
	ServiceLinkedIn:  true,
}

// VerificationService handles extraction of verification codes from emails.
//
// This service uses pattern matching to find and extract verification codes
// sent by popular services. It scans the inbox for the most recent matching
// email and returns the extracted code.
type VerificationService struct {
	client *Client
}

// GetCode extracts a verification code from the most recent email matching
// the specified service.
//
// Supported services: facebook, twitter, google, icloud, instagram,
// tiktok, discord, linkedin. You may also use the Service* constants.
//
//	code, err := client.Verification.GetCode(ctx, evilmail.ServiceGoogle, "user@evilmail.pro")
//	fmt.Printf("Your code is: %s\n", code.Code)
func (s *VerificationService) GetCode(ctx context.Context, service, email string) (*VerificationCode, error) {
	if service == "" {
		return nil, &ValidationError{Field: "service", Message: "service name is required"}
	}
	if !validServices[service] {
		return nil, &ValidationError{
			Field:   "service",
			Message: fmt.Sprintf("unsupported service %q; supported: facebook, twitter, google, icloud, instagram, tiktok, discord, linkedin", service),
		}
	}
	if email == "" {
		return nil, &ValidationError{Field: "email", Message: "email is required"}
	}

	path := fmt.Sprintf("/api/regex/%s?email=%s", url.PathEscape(service), url.QueryEscape(email))
	result, err := getJSON[VerificationCode](s.client, ctx, path)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
