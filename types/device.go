package types

type KeyloggerInputAllOS struct {
	// linux
	UsbName string `json:"usb_name"`

	// macOS
	VendorID  Hex `json:"vendor_id"`
	ProductID Hex `json:"product_id"`
}
