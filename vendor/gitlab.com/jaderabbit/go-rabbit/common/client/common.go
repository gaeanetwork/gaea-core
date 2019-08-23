package client

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/core/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func createClientConn(address, tlsRootCertFile string) (*grpc.ClientConn, error) {
	addr, override, clientConfig, err := configFromEnv("peer")
	if err != nil {
		return nil, err
	}

	if address == "" {
		address = addr
	}

	if clientConfig.SecOpts.UseTLS {
		if tlsRootCertFile == "" {
			return nil, errors.New("tls root cert file must be set")
		}
		caPEM, res := ioutil.ReadFile(tlsRootCertFile)
		if res != nil {
			err = errors.WithMessage(res, fmt.Sprintf("Unable to load TLS root cert file from %s", tlsRootCertFile))
			return nil, err
		}
		clientConfig.SecOpts.ServerRootCAs = [][]byte{caPEM}
	}

	gClient, err := comm.NewGRPCClient(clientConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create GRPCClient from clientConfig")
	}

	return gClient.NewConnection(address, override)
}

func configFromEnv(prefix string) (address, override string, clientConfig comm.ClientConfig, err error) {
	address = viper.GetString(prefix + ".address")
	override = viper.GetString(prefix + ".tls.serverhostoverride")
	clientConfig = comm.ClientConfig{}
	connTimeout := viper.GetDuration(prefix + ".client.connTimeout")
	if connTimeout == time.Duration(0) {
		connTimeout = defaultConnTimeout
	}
	clientConfig.Timeout = connTimeout
	secOpts := &comm.SecureOptions{
		UseTLS:            viper.GetBool(prefix + ".tls.enabled"),
		RequireClientCert: viper.GetBool(prefix + ".tls.clientAuthRequired")}
	if secOpts.UseTLS {
		caPEM, res := ioutil.ReadFile(config.GetPath(prefix + ".tls.rootcert.file"))
		if res != nil {
			err = errors.WithMessage(res,
				fmt.Sprintf("unable to load %s.tls.rootcert.file", prefix))
			return
		}
		secOpts.ServerRootCAs = [][]byte{caPEM}
	}
	if secOpts.RequireClientCert {
		keyPEM, res := ioutil.ReadFile(config.GetPath(prefix + ".tls.clientKey.file"))
		if res != nil {
			err = errors.WithMessage(res,
				fmt.Sprintf("unable to load %s.tls.clientKey.file", prefix))
			return
		}
		secOpts.Key = keyPEM
		certPEM, res := ioutil.ReadFile(config.GetPath(prefix + ".tls.clientCert.file"))
		if res != nil {
			err = errors.WithMessage(res,
				fmt.Sprintf("unable to load %s.tls.clientCert.file", prefix))
			return
		}
		secOpts.Certificate = certPEM
	}
	clientConfig.SecOpts = secOpts
	return
}
