package model

// Vendor represents a company being analyzed within a market.
// A vendor may have products in multiple product categories.
type Vendor struct {
	// ID is a unique identifier for the vendor.
	ID string `json:"id"`

	// Name is the company name.
	Name string `json:"name"`

	// Description provides context about the vendor.
	Description string `json:"description,omitempty"`

	// Color is the hex color code for this vendor in visualizations.
	// Example: "#8b5cf6" for purple. If not set, a default will be assigned.
	Color string `json:"color,omitempty"`

	// Website is the vendor's primary website.
	Website string `json:"website,omitempty"`

	// TargetSegments indicates which segments this vendor primarily serves.
	// Useful for understanding vendor positioning.
	TargetSegments []string `json:"targetSegments,omitempty"`

	// Founded is the year the company was founded (optional).
	Founded int `json:"founded,omitempty"`

	// Headquarters is the company's HQ location (optional).
	Headquarters string `json:"headquarters,omitempty"`

	// IsPublic indicates whether the company is publicly traded.
	IsPublic bool `json:"isPublic,omitempty"`

	// Ticker is the stock ticker symbol if public.
	Ticker string `json:"ticker,omitempty"`
}

// VendorProduct represents a vendor's offering in a specific product category.
// This allows tracking which vendors compete in which categories.
type VendorProduct struct {
	// VendorID references the vendor.
	VendorID string `json:"vendorId"`

	// CategoryID references the product category.
	CategoryID string `json:"categoryId"`

	// ProductName is the specific product name for this offering.
	ProductName string `json:"productName,omitempty"`

	// Description provides context about this product offering.
	Description string `json:"description,omitempty"`

	// LaunchYear is when this product was launched (optional).
	LaunchYear int `json:"launchYear,omitempty"`

	// IsPrimary indicates if this is the vendor's primary/flagship product.
	IsPrimary bool `json:"isPrimary,omitempty"`
}
