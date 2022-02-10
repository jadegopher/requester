package getter

import (
	"context"
	"io"
	"net/http"
	"time"
)

type HashGetter interface {
	GetResponseHash(ctx context.Context, to string) (hash []byte, err error)
}

type Getter struct {
	timeout time.Duration
	hash    iHash
}

type iHash interface {
	Calculate([]byte) []byte
}

func NewGetter(timeout time.Duration, hash iHash) *Getter {
	return &Getter{
		timeout: timeout,
		hash:    hash,
	}
}

func (g *Getter) GetResponseHash(ctx context.Context, to string) (hash []byte, err error) {
	resp, err := g.sendRequest(ctx, to)
	if err != nil {
		return nil, err
	}

	hash = g.hash.Calculate(resp)

	return hash, nil
}

func (g *Getter) sendRequest(ctx context.Context, to string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, to, nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
