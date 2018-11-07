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

type IdentityVerificationStatus struct {
	// Identity verification status code:<br> - VERIFIED: the identity is verified: the identity URL succeeded to sign a secret using the proof receipt's `pubKey` public key<br> - HTTP_ERROR: the identity URL returned an HTTP error<br> - INVALID_SIGNATURE: the identity URL returned an invalid signature (and thus failed to prove the ownership of the proof receipt's `pubKey` public key) 
	Code string `json:"code,omitempty"`
	// Identity verification status text giving more insight about verification errors.
	Text string `json:"text,omitempty"`
	// Array of X500 subject and issuer distinguished names of all X509 certificates of the identity URL.
	Certificates []X509SubjectIssuer `json:"certificates,omitempty"`
}