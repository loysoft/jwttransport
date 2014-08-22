package jwttransport

// Cache specifies the methods that implement a Token cache.
type Cache interface {
	Token() (*Token, error)
	PutToken(*Token) error
}

type CacheInMemory struct {
	token *Token
}

func (m CacheInMemory) Token() (*Token, error) {
	return m.token, nil
}

func (m CacheInMemory) PutToken(tok *Token) error {
	m.token = tok
	return nil
}
