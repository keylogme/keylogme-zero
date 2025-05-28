package types

func (d *KeyloggerInputAllOS) GetDeviceInput() KeyloggerInput {
	return KeyloggerInput{
		UsbName: d.UsbName,
	}
}

type KeyloggerInput struct {
	UsbName string `json:"usb_name"`
}
