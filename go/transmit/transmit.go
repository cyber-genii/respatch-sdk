// Package transmit is a deprecated alias for the respatch package.
//
// The SDK was renamed when the platform became Respatch. This package forwards every
// exported name to github.com/respatch/respatch-go-sdk/respatch so existing code keeps
// compiling; migrate by changing the import path and package selector:
//
//	-	import "github.com/transmit/transmit-go-sdk/transmit"
//	-	client, err := transmit.NewClient(transmit.Options{APIKey: key})
//	+	import "github.com/respatch/respatch-go-sdk/respatch"
//	+	client, err := respatch.NewClient(respatch.Options{APIKey: key})
//
// Because these are type aliases rather than wrapper types, values are
// interchangeable between the two packages during a gradual migration.
//
// Deprecated: use github.com/respatch/respatch-go-sdk/respatch.
package transmit

import "github.com/respatch/respatch-go-sdk/respatch"

// Client and configuration.
type (
	// Deprecated: use respatch.Client.
	Client = respatch.Client
	// Deprecated: use respatch.Options.
	Options = respatch.Options
)

// Services.
type (
	// Deprecated: use respatch.OrdersService.
	OrdersService = respatch.OrdersService
	// Deprecated: use respatch.DeliveriesService.
	DeliveriesService = respatch.DeliveriesService
	// Deprecated: use respatch.APIKeysService.
	APIKeysService = respatch.APIKeysService
	// Deprecated: use respatch.WebhooksService.
	WebhooksService = respatch.WebhooksService
	// Deprecated: use respatch.SandboxService.
	SandboxService = respatch.SandboxService
	// Deprecated: use respatch.VehiclesService.
	VehiclesService = respatch.VehiclesService
)

// Request and payload types.
type (
	// Deprecated: use respatch.QuoteLocationInput.
	QuoteLocationInput = respatch.QuoteLocationInput
	// Deprecated: use respatch.PackageDimensions.
	PackageDimensions = respatch.PackageDimensions
	// Deprecated: use respatch.QuotePackageInput.
	QuotePackageInput = respatch.QuotePackageInput
	// Deprecated: use respatch.GetQuoteRequest.
	GetQuoteRequest = respatch.GetQuoteRequest
	// Deprecated: use respatch.OrderContactLocation.
	OrderContactLocation = respatch.OrderContactLocation
	// Deprecated: use respatch.SupplierInfo.
	SupplierInfo = respatch.SupplierInfo
	// Deprecated: use respatch.OrderPackageMetaInput.
	OrderPackageMetaInput = respatch.OrderPackageMetaInput
	// Deprecated: use respatch.CreateDeliveryOrderRequest.
	CreateDeliveryOrderRequest = respatch.CreateDeliveryOrderRequest
	// Deprecated: use respatch.CalculateOrderFareRequest.
	CalculateOrderFareRequest = respatch.CalculateOrderFareRequest
	// Deprecated: use respatch.APIDeliveryRequest.
	APIDeliveryRequest = respatch.APIDeliveryRequest
	// Deprecated: use respatch.FareQuoteRequest.
	FareQuoteRequest = respatch.FareQuoteRequest
	// Deprecated: use respatch.CreateWebhookRequest.
	CreateWebhookRequest = respatch.CreateWebhookRequest
)

// NewClient creates a new Respatch API client.
//
// Deprecated: use respatch.NewClient.
func NewClient(options Options) (*Client, error) {
	return respatch.NewClient(options)
}

// VerifyWebhookSignature reports whether a webhook payload matches its signature.
//
// Deprecated: use respatch.VerifyWebhookSignature.
func VerifyWebhookSignature(payload []byte, signatureHeader, secret string) bool {
	return respatch.VerifyWebhookSignature(payload, signatureHeader, secret)
}

// SignWebhookPayload produces the signature header value for a payload.
//
// Deprecated: use respatch.SignWebhookPayload.
func SignWebhookPayload(payload []byte, secret string) string {
	return respatch.SignWebhookPayload(payload, secret)
}
