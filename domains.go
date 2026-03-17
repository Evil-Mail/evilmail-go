package evilmail

import "context"

// DomainsService handles operations related to available email domains.
// This endpoint requires authentication and returns customer-specific domains.
type DomainsService struct {
	client *Client
}

// List retrieves all available email domains, grouped by type.
// The response includes free, premium, and customer-owned domains,
// along with the authenticated user's package type.
//
//	domains, err := client.Domains.List(ctx)
//	fmt.Println("Free domains:", domains.Free)
//	fmt.Println("Premium domains:", domains.Premium)
//	fmt.Println("Your package:", domains.PackageType)
func (s *DomainsService) List(ctx context.Context) (*Domains, error) {
	result, err := getJSON[Domains](s.client, ctx, "/api/ext/customer-domains")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PublicDomainsService handles the public (unauthenticated) domains listing.
type PublicDomainsService struct {
	client *Client
}

// List retrieves publicly available email domains.
// This endpoint does not require authentication.
//
//	domains, err := client.PublicDomains.List(ctx)
//	fmt.Println("Free domains:", domains.Free)
//	fmt.Println("Premium domains:", domains.Premium)
func (s *PublicDomainsService) List(ctx context.Context) (*PublicDomains, error) {
	result, err := getJSON[PublicDomains](s.client, ctx, "/api/ext/domains")
	if err != nil {
		return nil, err
	}
	return &result, nil
}
