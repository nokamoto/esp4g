package esp4g

type proxyMessage struct {
	bytes []byte
}

func newProxyMessage() *proxyMessage {
	return &proxyMessage{}
}

func (p *proxyMessage)Marshal() ([]byte, error) {
	return p.bytes, nil
}

func (p *proxyMessage)Unmarshal(bytes []byte) error {
	p.bytes = bytes
	return nil
}

func (p *proxyMessage)Reset() {
	*p = proxyMessage{}
}

func (p *proxyMessage)String() string {
	return string(p.bytes)
}

func (*proxyMessage)ProtoMessage() {}
