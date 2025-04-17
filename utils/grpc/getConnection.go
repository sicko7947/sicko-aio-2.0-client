package grpc

import "github.com/shimingyah/pool"

func GetFunctionServerGrpcCon() (pool.Conn, error) {
	return functionServerGRPCPool.Get()
}

func GetAuthServerGrpcCon() (pool.Conn, error) {
	return authServerGRPCPool.Get()
}
