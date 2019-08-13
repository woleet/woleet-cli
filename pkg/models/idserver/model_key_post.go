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

type KeyPost struct {
	// Key name.
	Name string `json:"name"`
	// Key expiration date (Unix ms timestamp). <br>Note that the field is not returned if the key has no expiration date. 
	Expiration int64 `json:"expiration,omitempty"`
	Status KeyStatusEnum `json:"status,omitempty"`
}
