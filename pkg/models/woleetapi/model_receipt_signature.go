/*
 * Woleet API
 *
 * Welcome to **Woleet API reference documentation**.<br> It is highly recommanded to read the chapters **[introducing Woleet API concepts](https://doc.woleet.io/v1.5.1/reference)** before reading this documentation. 
 *
 * API version: 1.5.1
 * Contact: contact@woleet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package woleetapi

type ReceiptSignature struct {
	SignedHash string `json:"signedHash"`
	PubKey string `json:"pubKey"`
	Signature string `json:"signature"`
	IdentityURL string `json:"identityURL,omitempty"`
}