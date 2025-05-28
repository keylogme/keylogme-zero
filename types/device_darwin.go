package types

func (d *KeyloggerInputAllOS) GetDeviceInput() KeyloggerInput {
	return KeyloggerInput{
		VendorID:  d.VendorID,
		ProductID: d.ProductID,
	}
}

type KeyloggerInput struct {
	VendorID  Hex `json:"vendor_id"`
	ProductID Hex `json:"product_id"`
}
