package grpc

import (
	"context"
)

// Token token认证
type Token struct {
	Value string
}

const headerAuthorize string = "authorization"

// GetRequestMetadata
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{headerAuthorize: t.Value}, nil
}

// RequireTransportSecurity
func (t *Token) RequireTransportSecurity() bool {
	return true
}
