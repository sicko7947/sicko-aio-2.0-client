package grpc

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"log"

	"github.com/shimingyah/pool"
	"github.com/sicko7947/sickocommon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/constants"
)

const (
	address     = "api.sickoaio.com:26501"
	authAddress = "auth.sickoaio.com:26501"
)

var (
	authServerGRPCPool     pool.Pool
	functionServerGRPCPool pool.Pool
	token                  Token
)

func init() {
	go ResetConn()
	functionServerGRPCPool = newFunctionServerConPool()
}

func ResetConn() {
	authServerGRPCPool = nil
	for {
		if len(communicator.Config.Settings.Key) > 0 {
			value, _ := sickocommon.RsaEncrypt([]byte(communicator.Config.Settings.Key), []byte(constants.AUTH_PUBLIC_KEY))
			vBase64 := base64.StdEncoding.EncodeToString(value)
			token = Token{
				Value: vBase64,
			}
			authServerGRPCPool = newAuthServerConPool()
			return
		}
	}
}

// newFunctionServerConPool : Create grpc Pool recycle connections
func newFunctionServerConPool() pool.Pool {
	grpcPool, err := pool.New(address, pool.Options{
		MaxIdle:              pool.DefaultOptions.MaxIdle,
		MaxActive:            pool.DefaultOptions.MaxActive,
		MaxConcurrentStreams: pool.DefaultOptions.MaxConcurrentStreams,
		Reuse:                pool.DefaultOptions.Reuse,
		Dial: func(address string) (*grpc.ClientConn, error) {
			// set credentials
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(constants.FUNCTION_SERVER_CERT))
			creds := credentials.NewClientTLSFromCert(caCertPool, `api.sickoaio.com`)

			ctx, cancel := context.WithTimeout(context.Background(), pool.DialTimeout)
			defer cancel()
			return grpc.DialContext(ctx, address,
				grpc.WithTransportCredentials(creds),
				// grpc.WithPerRPCCredentials(&token),
				grpc.WithInitialWindowSize(pool.InitialWindowSize),
				grpc.WithInitialConnWindowSize(pool.InitialConnWindowSize),
				grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(pool.MaxSendMsgSize)),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(pool.MaxRecvMsgSize)),
				grpc.WithKeepaliveParams(keepalive.ClientParameters{
					Time:                pool.KeepAliveTime,
					Timeout:             pool.KeepAliveTimeout,
					PermitWithoutStream: true,
				}),
			)
		},
	})
	if err != nil {
		log.Fatalf("Failed to create connection pool %v", err)
	}
	return grpcPool
}

// newAuthServerConPool : Create grpc Pool recycle connections
func newAuthServerConPool() pool.Pool {
	grpcPool, err := pool.New(authAddress, pool.Options{
		MaxIdle:              pool.DefaultOptions.MaxIdle,
		MaxActive:            pool.DefaultOptions.MaxActive,
		MaxConcurrentStreams: pool.DefaultOptions.MaxConcurrentStreams,
		Reuse:                pool.DefaultOptions.Reuse,
		Dial: func(address string) (*grpc.ClientConn, error) {
			// set credentials
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(constants.AUTH_SERVER_CERT))
			creds := credentials.NewClientTLSFromCert(caCertPool, `auth.sickoaio.com`)

			ctx, cancel := context.WithTimeout(context.Background(), pool.DialTimeout)
			defer cancel()
			return grpc.DialContext(ctx, address,
				grpc.WithTransportCredentials(creds),
				grpc.WithPerRPCCredentials(&token),
				grpc.WithInitialWindowSize(pool.InitialWindowSize),
				grpc.WithInitialConnWindowSize(pool.InitialConnWindowSize),
				grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(pool.MaxSendMsgSize)),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(pool.MaxRecvMsgSize)),
				grpc.WithKeepaliveParams(keepalive.ClientParameters{
					Time:                pool.KeepAliveTime,
					Timeout:             pool.KeepAliveTimeout,
					PermitWithoutStream: true,
				}),
			)
		},
	})
	if err != nil {
		log.Fatalf("Failed to create connection pool %v", err)
	}
	return grpcPool
}
