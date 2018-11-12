/*
 * Woleet.ID Server
 *
 * This is Woleet.ID Server API documentation.
 *
 * API version: 1.0.4
 * Contact: contact@woleet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package idserver

type SignatureResult struct {
	// The public key used to sign (must be the same as the `pubKey` parameter, if provided).
	PubKey string `json:"pubKey,omitempty"`
	// The hash that is signed (same as the `hashToSign` parameter).
	SignedHash string `json:"signedHash,omitempty"`
	// The signature of `hashToSign` using the public key `pubKey`.
	Signature string `json:"signature,omitempty"`
	// The public URL of the `/identity` endpoint (ie. a URL that anyone can use to prove and verify the identity associated with the public key).
	IdentityURL string `json:"identityURL,omitempty"`
}
