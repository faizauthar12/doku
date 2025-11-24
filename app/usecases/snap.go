package usecases

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

func readPrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	var rsaPrivateKey *rsa.PrivateKey

	if block.Type == "PRIVATE KEY" {
		pkcs8PrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
		}

		pkcs8RSAPrivateKey, ok := pkcs8PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}

		rsaPrivateKey = pkcs8RSAPrivateKey
	} else if block.Type == "RSA PRIVATE KEY" {
		pkcs1PrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS1 private key: %w", err)
		}

		rsaPrivateKey = pkcs1PrivateKey
	} else {
		return nil, errors.New("unsupported private key type")
	}

	return rsaPrivateKey, nil
}

func generateTokenStringToSign(clientId, xTimestamp string) string {
	return fmt.Sprintf("%s|%s", clientId, xTimestamp)
}

func generateGetTokenSignature(privateKey, xTimestamp, clientId string) (string, error) {

	rsaPrivateKey, err := readPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	fmt.Printf("Using RSA Private Key : %v\n", rsaPrivateKey)
	stringToSign := generateTokenStringToSign(clientId, xTimestamp)
	stringToSignHash := sha256.Sum256([]byte(stringToSign))

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, stringToSignHash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func generateRequestStringToSign(httpMethod, requestPath, accessToken, xTimestamp string, jsonBody string) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", httpMethod, requestPath, accessToken, jsonBody, xTimestamp)
}

func generateKirimDokuRequestSignature(clientSecret, httpMethod, requestPath, accessToken, xTimestamp string, jsonBody []byte) (string, error) {
	// rsaPrivateKey, err := readPrivateKey(privateKey)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to read private key: %w", err)
	// }

	encodedRequestBody := sha256.Sum256(jsonBody)
	hexEncodedRequestBody := hex.EncodeToString(encodedRequestBody[:])
	hexEncodedRequestBody = strings.ToLower(hexEncodedRequestBody)

	stringToSign := generateRequestStringToSign(httpMethod, requestPath, accessToken, xTimestamp, hexEncodedRequestBody)
	fmt.Println("String to Sign:", stringToSign)

	hmac := hmac.New(sha512.New, []byte(clientSecret))
	hmac.Write([]byte(stringToSign))
	hmacSum := hmac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(hmacSum)
	fmt.Println("Signature:", signature)

	return signature, nil
}
