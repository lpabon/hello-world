package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"time"

	"github.com/libopenstorage/openstorage/pkg/grpcserver"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	endpoint string
	port     int
)

type portworxGrpcConnection struct {
	conn        *grpc.ClientConn
	dialOptions []grpc.DialOption
	endpoint    string
}

func (pg *portworxGrpcConnection) setDialOptions(tls bool) error {
	if tls {
		// Setup a connection
		capool, err := x509.SystemCertPool()
		if err != nil {
			return fmt.Errorf("Failed to load CA system certs: %v", err)
		}
		pg.dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(capool, ""),
		)}
	} else {
		pg.dialOptions = []grpc.DialOption{grpc.WithInsecure()}
	}

	return nil
}

func (pg *portworxGrpcConnection) getGrpcConn() (*grpc.ClientConn, error) {

	if pg.conn == nil {
		var err error
		pg.conn, err = grpcserver.Connect(pg.endpoint, pg.dialOptions)
		if err != nil {
			return nil, fmt.Errorf("Error connecting to GRPC server[%s]: %v", pg.endpoint, err)
		}
	}
	return pg.conn, nil
}

func main() {

	flag.StringVar(&endpoint, "endpoint", "", "Endpoint to use")
	flag.IntVar(&port, "port", 9020, "Port to use")

	flag.Parse()

	if len(endpoint) == 0 {
		logrus.Infof("endpoint is required")
		return
	}

	go func() {
		sdkConn := &portworxGrpcConnection{
			endpoint: fmt.Sprintf("%s:%d", endpoint, port),
		}
		err := sdkConn.setDialOptions(false)
		if err != nil {
			logrus.Infof("Error: setDialOptions due to: %v", err)
			return
		}

		conn, err := sdkConn.getGrpcConn()
		if err != nil {
			logrus.Infof("Error: getGrpcConn due to: %v", err)
			return
		}

		logrus.Infof("Successfully got connection: %v", conn)
	}()

	time.Sleep(5 * time.Hour)

}
