package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"os"
	"time"
)

func main() {
	err := genToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := genKey(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func genToken() error {
	file, err := os.Open("deployments/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem")
	if err != nil {
		return err
	}

	// limit PEM file size to 1 megabyte. This should bew reasonable for
	// almost any PEM file and prevents undefined behaviour like linking the file
	// to /dev/random or something like that
	privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return fmt.Errorf("error reading auth private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("error parsing auth private key: %w", err)
	}
	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.

	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to
	// be used only once)

	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   "123456789", //TO BE REPLACED
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	method := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(method, claims)
	// Adding a new section to the header for verifying the private key ID
	token.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	str, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("error signing token: %w", err)
	}

	fmt.Println("========= TOKEN BEGIN =========")
	fmt.Println(str)
	fmt.Println("========= TOKEN END =========")
	fmt.Println()

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("error marshaling public key: %w", err)
	}

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the private key file.
	if err := pem.Encode(os.Stdout, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	return nil
}

// genKey creates a x509 private/public key for auth tokens.
func genKey() error {

	// Generate a new private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create a file for the private key information in PEM form.
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("error creating private file: %w", err)
	}
	defer privateFile.Close()

	// Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("error encoding to private file: %w", err)
	}

	// ===============================================================

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("error marshaling public key: %w", err)
	}

	// Create a file for the public key information in PEM form.
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("error creating public file: %w", err)
	}
	defer publicFile.Close()

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the private key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	fmt.Println("private and public key files generated")

	return nil
}
