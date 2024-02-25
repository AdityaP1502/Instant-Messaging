package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/youmark/pkcs8"
)

type envKey string

var (
	PASSPHRASE       envKey = "PASSPHRASE"
	CERT_FILE_PATH   envKey = "CERT_FILE_PATH"
	PRIVATE_KEY_PATH envKey = "PRIVATE_KEY_PATH"
	HOST             envKey = "HOST"
)

func SignCertificate(csrFile io.Reader, caCRT *x509.Certificate, caPrivateKey interface{}) (*os.File, error) {
	// load client certificate request
	clientCSRFile, err := io.ReadAll(csrFile)

	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(clientCSRFile)
	if pemBlock == nil {
		return nil, err
	}
	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	if err = clientCSR.CheckSignature(); err != nil {
		return nil, err
	}

	// create client certificate template
	clientCRTTemplate := x509.Certificate{
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,

		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		PublicKey:          clientCSR.PublicKey,

		SerialNumber: big.NewInt(2),
		Issuer:       caCRT.Subject,
		Subject:      clientCSR.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// create client certificate from template and CA public key
	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCRT, clientCSR.PublicKey, caPrivateKey)
	if err != nil {
		return nil, err
	}

	// save the certificate
	clientCRTFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return nil, err
	}
	pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: clientCRTRaw})

	return clientCRTFile, nil
}

func main() {
	pw := []byte(os.Getenv(string(PASSPHRASE)))

	caPublicKeyFile, err := os.ReadFile(os.Getenv(string(CERT_FILE_PATH)))

	if err != nil {
		log.Fatal(err)
	}

	pemBlock, _ := pem.Decode(caPublicKeyFile)

	if pemBlock == nil {
		log.Fatal("pem decode failed")
	}

	caCert, err := x509.ParseCertificate(
		pemBlock.Bytes,
	)

	if err != nil {
		log.Fatal(err)
	}

	caPrivateKeyFile, err := os.ReadFile(os.Getenv(string(PRIVATE_KEY_PATH)))
	if err != nil {
		log.Fatal(err)
	}

	pemBlock, _ = pem.Decode(caPrivateKeyFile)

	if pemBlock == nil {
		log.Fatal("pem decode failed")
	}

	caPrivateKey, err := pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, pw)

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/certificate/sign", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20) // 10 mb
		csrFile, _, err := r.FormFile("csr")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer csrFile.Close()

		crtFile, err := SignCertificate(csrFile, caCert, caPrivateKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, crtFile.Name())

	}).Methods("POST")

	var wg sync.WaitGroup

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{caCert.Raw},
				PrivateKey:  caPrivateKey,
			},
		},
		MinVersion: tls.VersionTLS11,
	}

	srv := http.Server{
		Addr:      os.Getenv(string(HOST)),
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	wg.Add(1)
	go func() {
		log.Fatal(srv.ListenAndServeTLS("", ""))
	}()

	fmt.Printf("Cert manager is running on %s\n", os.Getenv(string(HOST)))
	wg.Wait()
}
