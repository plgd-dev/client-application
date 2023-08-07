package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/apex/log"
	externalip "github.com/glendc/go-external-ip"
)

func getExternalIP() (net.IP, error) {
	consensus := externalip.DefaultConsensus(&externalip.ConsensusConfig{
		Timeout: time.Second * 3,
	}, nil)
	return consensus.ExternalIP()
}

func getExternalDomains(ip net.IP) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	dns, err := net.DefaultResolver.LookupAddr(ctx, ip.String())
	if err != nil {
		return nil, err
	}
	// remove trailing dot
	for i := range dns {
		dns[i] = strings.TrimSuffix(dns[i], ".")
	}
	return dns, nil
}

func setupIPsFromInterface(template *x509.Certificate, iface net.Interface) {
	if iface.Flags&net.FlagLoopback != 0 {
		return
	}
	if iface.Flags&net.FlagUp == 0 {
		return
	}
	addrs, err := iface.Addrs()
	if err != nil {
		log.Warnf("cannot get network interface(%v) addresses: %v", iface.Name, err)
		return
	}
	for _, addr := range addrs {
		if ip, _, err3 := net.ParseCIDR(addr.String()); err3 == nil && !ip.IsUnspecified() {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			log.Warnf("cannot parse network interface(%v) address(%v): %v", iface.Name, addr.String(), err3)
		}
	}
}

func setupIPsFromInterfaces(template *x509.Certificate) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Warnf("cannot get network interfaces: %v", err)
		return
	}
	for _, iface := range ifaces {
		setupIPsFromInterface(template, iface)
	}
}

func setupIPAndDomains(template *x509.Certificate) {
	template.DNSNames = []string{"localhost"}
	template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}
	if ip, err := getExternalIP(); err != nil {
		log.Warnf("cannot get external IP: %v", err)
	} else {
		template.IPAddresses = append(template.IPAddresses, ip)
		domains, err := getExternalDomains(ip)
		if err != nil {
			log.Warnf("cannot get external domains: %v", err)
		} else {
			template.DNSNames = append(template.DNSNames, domains...)
		}
	}
	setupIPsFromInterfaces(template)
}

const cn = "self-signed-certificate"

func checkSelfSignedCertificate(certFile, keyFile string) bool {
	crt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return false
	}
	cert, err := x509.ParseCertificate(crt.Certificate[0])
	if err != nil {
		return false
	}
	if cert.Subject.CommonName != cn {
		// it is not self signed certificate
		return true
	}
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return false
	}
	return true
}

func generateSelfSigned(certFile, keyFile string) error {
	err := os.MkdirAll(path.Dir(certFile), 0o700)
	if err != nil {
		return fmt.Errorf("cannot create directory(%v) for certificate: %w", path.Dir(certFile), err)
	}
	err = os.MkdirAll(path.Dir(keyFile), 0o700)
	if err != nil {
		return fmt.Errorf("cannot create directory(%v) for key: %w", path.Dir(keyFile), err)
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("cannot generate key: %w", err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"plgd.dev"},
			CommonName:   cn,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	setupIPAndDomains(&template)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("cannot create certificate: %w", err)
	}

	derKeyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("cannot marshal private key: %w", err)
	}

	err = os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}), 0o600)
	if err != nil {
		return fmt.Errorf("cannot write certificate to file(%v): %w", certFile, err)
	}

	err = os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derKeyBytes,
	}), 0o600)
	if err != nil {
		return fmt.Errorf("cannot write key to file(%v): %w", keyFile, err)
	}

	return nil
}
