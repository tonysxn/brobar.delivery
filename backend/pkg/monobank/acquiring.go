package monobank

type Acquiring struct {
	xToken       string
	publicDomain string
}

func NewAcquiring(xToken string, publicDomain string) *Acquiring {
	return &Acquiring{
		xToken:       xToken,
		publicDomain: publicDomain,
	}
}
