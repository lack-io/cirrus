package client

type Option struct {
	Headless bool

	BlinkSettings string

	UserAgent string

	IgnoreCertificateErrors bool

	WindowsHigh int

	WindowsWith int
}
