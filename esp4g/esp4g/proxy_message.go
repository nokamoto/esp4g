package esp4g

type ProxyMessage struct {
	bytes []byte
}

func NewProxyMessage() *ProxyMessage {
	return &ProxyMessage{}
}

func (p *ProxyMessage)Marshal() ([]byte, error) {
	return p.bytes, nil
}

func (p *ProxyMessage)Unmarshal(bytes []byte) error {
	p.bytes = bytes
	return nil
}

func (p *ProxyMessage)Reset() {
	*p = ProxyMessage{}
}

func (p *ProxyMessage)String() string {
	return string(p.bytes)
}

func (*ProxyMessage)ProtoMessage() {}
