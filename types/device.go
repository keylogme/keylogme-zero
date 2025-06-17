package types

type KeyloggerInput struct {
	VendorID  Hex `json:"vendor_id"`
	ProductID Hex `json:"product_id"`
}
