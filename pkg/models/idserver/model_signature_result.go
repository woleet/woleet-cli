/*
 * Woleet.ID Server
 *
 * This is Woleet.ID Server API documentation.
 *
 * API version: 1.2.5
 * Contact: contact@woleet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package idserver
// SignatureResult struct for SignatureResult
type SignatureResult struct {
	// Public key used to sign (must be the same as the `pubKey` parameter if provided).
	PubKey string `json:"pubKey,omitempty"`
	// SHA256 hash that is signed (same as the `hashToSign` parameter).
	SignedHash string `json:"signedHash,omitempty"`
	// Message that is signed (same as the `messageToSign` parameter).
	SignedMessage string `json:"signedMessage,omitempty"`
	// Signature of `messageToSign` or `hashToSign` using the public key `pubKey`, or signature of SHA256(`signedMessage` or `signedHash` + `signedIdentity` + `signedIssuerDomain`) if the identity of the signer and the domain of the identity issuer are included to the signed data. 
	Signature string `json:"signature,omitempty"`
	// Public URL of the **Identity endpoint** (ie. the URL that anyone can use to get the identity associated to a public key). 
	IdentityURL string `json:"identityURL,omitempty"`
	// X500 Distinguished Name representing the identity of the signer.<br> Returned only if `identityToSign` is used. 
	SignedIdentity string `json:"signedIdentity,omitempty"`
	// Domain of the identity issuer (ie. of the organization who verified the identity).<br> Returned only if `identityToSign` is used. 
	SignedIssuerDomain string `json:"signedIssuerDomain,omitempty"`
}
