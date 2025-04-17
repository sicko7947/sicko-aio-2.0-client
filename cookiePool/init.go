package cookiePool

import (
	"github.com/shimingyah/pool"
	auth_service "sicko-aio-2.0-client/proto/auth"
)

var (
	conn             pool.Conn
	cookieDataStream auth_service.AuthStream_RequestCookieDataClient
)

// func init() {
// 	t := time.NewTicker(3 * time.Second)
// 	go func() {
// 		for {
// 			if conn == nil || conn.Value().GetState().String() != "READY" {
// 				conn, err := grpc.GrpcConnPool.Get()
// 				if err != nil {
// 					conn = nil
// 					continue
// 				}
// 				streamClient := grpc_service.NewStreamClient(conn.Value())
// 				if cookieDataStream, err = streamClient.RequestCookieData(context.Background()); err != nil {
// 					conn = nil
// 					continue
// 				}
// 			}
// 			<-t.C
// 		}
// 	}()
// }
