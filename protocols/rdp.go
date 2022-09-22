package protocols

// NewRDPProtocol initializes a Protocol with a SOCKS4 signature.
func NewRDPProtocol(targetAddress string) *Protocol {
	return &Protocol{
		Name:            "RDP",
		Target:          targetAddress,
		MatchStartBytes: [][]byte{{0x03, 0x00, 0x00}},
	}
}
