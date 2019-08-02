/*
 * Woleet.ID Server
 *
 * This is Woleet.ID Server API documentation.
 *
 * API version: 1.2.0
 * Contact: contact@woleet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package idserver

type IdentityResultAllOf struct {
	// The right part of the signed data (generated randomly). <br>To prevent man-in-the-middle attacks, the data starts with the server's identity URL and this should be verified by the caller. 
	RightData string `json:"rightData,omitempty"`
	// The signature of the concatenation of `leftData` and `rightData` using the public key `pubKey`. 
	Signature string `json:"signature,omitempty"`
	Identity Identity `json:"identity,omitempty"`
	Key Key `json:"key,omitempty"`
}
