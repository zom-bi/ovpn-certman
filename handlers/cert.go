package handlers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zom-bi/ovpn-certman/models"
	"github.com/zom-bi/ovpn-certman/services"
	"github.com/go-chi/chi"

	"github.com/zom-bi/ovpn-certman/views"
)

func ListClientsHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.NewWithSession(req, p.Sessions)

		username := p.Sessions.GetUsername(req)

		clients, _ := p.ClientCollection.ListClientsForUser(username)

		v.Vars["Clients"] = clients
		v.Render(w, "client_list")
	}
}

func CreateCertHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		username := p.Sessions.GetUsername(req)
		certname := req.FormValue("certname")

		// Validate certificate Name
		if !IsByteLength(certname, 2, 64) || !IsDNSName(certname) {
			p.Sessions.Flash(w, req,
				services.Flash{
					Type:    "danger",
					Message: "The certificate name can only contain letters and numbers",
				},
			)
			http.Redirect(w, req, "/certs", http.StatusFound)
			return
		}

		// lowercase the certificate name, to avoid problems with the case
		// insensitive matching inside OpenVPN
		certname = strings.ToLower(certname)

		// Load CA master certificate
		caCert, caKey, err := loadX509KeyPair("ca.crt", "ca.key")
		if err != nil {
			log.Fatalf("error loading ca keyfiles: %s", err)
			panic(err.Error())
		}

		// Generate Keypair
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("Could not generate keypair: %s", err)
			p.Sessions.Flash(w, req,
				services.Flash{
					Type:    "danger",
					Message: "The certificate key could not be generated",
				},
			)
			http.Redirect(w, req, "/certs", http.StatusFound)
			return
		}

		// Generate Certificate
		commonName := fmt.Sprintf("%s@%s", username, certname)
		derBytes, err := CreateCertificate(commonName, key, caCert, caKey)

		// Initialize new client config
		client := models.Client{
			Name:       certname,
			CreatedAt:  time.Now(),
			PrivateKey: x509.MarshalPKCS1PrivateKey(key),
			Cert:       derBytes,
			User:       username,
		}

		// Insert client into database
		if err := p.ClientCollection.CreateClient(&client); err != nil {
			log.Println(err.Error())
			p.Sessions.Flash(w, req,
				services.Flash{
					Type:    "danger",
					Message: "The certificate could not be added to the database",
				},
			)
			http.Redirect(w, req, "/certs", http.StatusFound)
			return
		}

		p.Sessions.Flash(w, req,
			services.Flash{
				Type:    "success",
				Message: "The certificate was created successfully",
			},
		)

		http.Redirect(w, req, "/certs", http.StatusFound)
	}
}

func DeleteCertHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.New(req)
		// detemine own username
		username := p.Sessions.GetUsername(req)
		name := chi.URLParam(req, "name")

		client, err := p.ClientCollection.GetClientByNameUser(name, username)
		if err != nil {
			v.RenderError(w, http.StatusNotFound)
			return
		}

		err = p.ClientCollection.DeleteClient(client.ID)
		if err != nil {
			p.Sessions.Flash(w, req,
				services.Flash{
					Type:    "danger",
					Message: "Failed to delete certificate",
				},
			)
			http.Redirect(w, req, "/certs", http.StatusFound)
		}

		p.Sessions.Flash(w, req,
			services.Flash{
				Type:    "success",
				Message: template.HTML(fmt.Sprintf("Successfully deleted client <strong>%s</strong>", client.Name)),
			},
		)
		http.Redirect(w, req, "/certs", http.StatusFound)
	}
}

func DownloadCertHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.New(req)
		// detemine own username
		username := p.Sessions.GetUsername(req)
		name := chi.URLParam(req, "name")

		client, err := p.ClientCollection.GetClientByNameUser(name, username)
		if err != nil {
			v.RenderError(w, http.StatusNotFound)
			return
		}

		// cbuf and kbuf are buffers in which the PEM certificates are
		// rendered into
		var cbuf = new(bytes.Buffer)
		var kbuf = new(bytes.Buffer)

		pem.Encode(cbuf, &pem.Block{Type: "CERTIFICATE", Bytes: client.Cert})
		pem.Encode(kbuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: client.PrivateKey})

		ca, err := ioutil.ReadFile("ca.crt")
		if err != nil {
			log.Printf("Error loading ca file: %s", err)
			v.RenderError(w, http.StatusInternalServerError)
			return
		}

		ta, err := ioutil.ReadFile("ta.key")
		if err != nil {
			log.Printf("Error loading ta file: %s", err)
			v.RenderError(w, http.StatusInternalServerError)
			return
		}

		vars := map[string]string{
			"CA":    string(ca),
			"TA":    string(ta),
			"Cert":  cbuf.String(),
			"Key":   kbuf.String(),
			"User":  username,
			"Name":  name,
			"Dev":   os.Getenv("VPN_DEV"),
			"Host":  os.Getenv("VPN_HOST"),
			"Port":  os.Getenv("VPN_PORT"),
			"Proto": os.Getenv("VPN_PROTO"),
		}

		t, err := views.GetTemplate("config.ovpn")
		if err != nil {
			log.Printf("Error loading certificate template: %s", err)
			v.RenderError(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/x-openvpn-profile")
		w.Header().Set("Content-Disposition", "attachment; filename=\"config.ovpn\"")
		w.WriteHeader(http.StatusOK)
		t.Execute(w, vars)
		return
	}
}

func loadX509KeyPair(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	cf, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}

	kf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	cpb, _ := pem.Decode(cf)
	kpb, _ := pem.Decode(kf)
	crt, err := x509.ParseCertificate(cpb.Bytes)

	if err != nil {
		return nil, nil, err
	}
	key, err := x509.ParsePKCS1PrivateKey(kpb.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return crt, key, nil
}

// CreateCertificate creates a CA-signed certificate
func CreateCertificate(commonName string, key interface{}, caCert *x509.Certificate, caKey interface{}) ([]byte, error) {
	subj := caCert.Subject
	// .. except for the common name
	subj.CommonName = commonName

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Obscure error in cert serial number generation: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subj,

		NotBefore: time.Now().Add(-5 * time.Minute),         // account for clock shift
		NotAfter:  time.Now().Add(24 * time.Hour * 356 * 5), // 5 years ought to be enough!

		SignatureAlgorithm:    x509.SHA256WithRSA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	return x509.CreateCertificate(rand.Reader, &template, caCert, publicKey(key), caKey)
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
