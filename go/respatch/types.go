package respatch

type QuoteLocationInput struct {
	Address *string  `json:"address,omitempty"`
	Lat     *float64 `json:"lat,omitempty"`
	Lng     *float64 `json:"lng,omitempty"`
}

type PackageDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type QuotePackageInput struct {
	WeightKg     float64            `json:"weight_kg"`
	DimensionsCm PackageDimensions  `json:"dimensions_cm"`
	Quantity     int                `json:"quantity,omitempty"`
}

type GetQuoteRequest struct {
	Pickup              QuoteLocationInput   `json:"pickup"`
	Dropoff             QuoteLocationInput   `json:"dropoff"`
	DeliveryType        string               `json:"delivery_type"`
	VehicleType         *string              `json:"vehicle_type,omitempty"`
	Packages            []QuotePackageInput  `json:"packages"`
	ScheduledPickupTime *string              `json:"scheduled_pickup_time,omitempty"`
}

type OrderContactLocation struct {
	ContactName           string  `json:"contact_name"`
	ContactEmail          *string `json:"contact_email,omitempty"`
	ContactPhone          string  `json:"contact_phone"`
	ContactPhoneSecondary *string `json:"contact_phone_secondary,omitempty"`
}

type SupplierInfo struct {
	Name      *string `json:"name,omitempty"`
	Reference *string `json:"reference,omitempty"`
}

type OrderPackageMetaInput struct {
	PackageID    *string       `json:"package_id,omitempty"`
	Description  string        `json:"description"`
	Value        float64       `json:"value,omitempty"`
	Fragile      bool          `json:"fragile,omitempty"`
	SupplierInfo *SupplierInfo `json:"supplier_info,omitempty"`
}

type CreateDeliveryOrderRequest struct {
	QuoteID           string                   `json:"quote_id"`
	Pickup            OrderContactLocation     `json:"pickup"`
	Dropoff           OrderContactLocation     `json:"dropoff"`
	Packages          []OrderPackageMetaInput  `json:"packages"`
	PaymentMethod     *string                  `json:"payment_method,omitempty"`
	WebhookURL        *string                  `json:"webhook_url,omitempty"`
	ExternalReference *string                  `json:"external_reference,omitempty"`
}

// Deprecated: use GetQuoteRequest
type CalculateOrderFareRequest = GetQuoteRequest

// Deprecated: use CreateDeliveryOrderRequest
type APIDeliveryRequest = CreateDeliveryOrderRequest

// Deprecated: use GetQuoteRequest
type FareQuoteRequest = GetQuoteRequest

type CreateWebhookRequest struct {
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Secret      string   `json:"secret,omitempty"`
	Description string   `json:"description,omitempty"`
}
