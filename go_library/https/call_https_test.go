package https

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestHttpsCall(t *testing.T) {
	file, err := ioutil.ReadFile("/home/jujucom/Desktop/workspace/projects/go_samples/go_library/https/ca.crt")
	if err != nil {
		log.Fatal(err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(file)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			//Certificates: []tls.Certificate{},
			RootCAs: caCertPool,
		},
	}

	client := &http.Client{Transport: transport}
	resp, err := client.Get("https://www.home.com")
	if err != nil {
		log.Fatal(err)
		return
	}
	htmlData, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	print(string(htmlData))
}
