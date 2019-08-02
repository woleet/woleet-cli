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

type KeyGetAllOf struct {
	// Key identifier (allocated by the platform).
	Id string `json:"id,omitempty"`
	// Public key (bitcoin address when using BIP39 keys).
	PubKey string `json:"pubKey,omitempty"`
	Type KeyTypeEnum `json:"type,omitempty"`
	Holder KeyHolderEnum `json:"holder,omitempty"`
	Device KeyDeviceEnum `json:"device,omitempty"`
	// Indicates whether the key has expired or not. <br>Note that the field is not returned if the key has not expired. 
	Expired *bool `json:"expired,omitempty"`
	// Date of creation (Unix ms timestamp).
	CreatedAt int64 `json:"createdAt,omitempty"`
	// Date of last modification (Unix ms timestamp).
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// Date of last usage (Unix ms timestamp).
	LastUsed int64 `json:"lastUsed,omitempty"`
}
