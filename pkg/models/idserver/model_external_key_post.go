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

type ExternalKeyPost struct {
	// Key name.
	Name string `json:"name"`
	// Key address.
	PublicKey string `json:"publicKey"`
	Device KeyDeviceEnum `json:"device,omitempty"`
	Status KeyStatusEnum `json:"status,omitempty"`
	// Key expiration date (Unix ms timestamp). 
	Expiration int64 `json:"expiration,omitempty"`
}