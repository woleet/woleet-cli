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

type UserGetAllOf struct {
	// Date of creation (Unix ms timestamp).
	CreatedAt int64 `json:"createdAt,omitempty"`
	// Date of last modification (Unix ms timestamp).
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// Date of last login (Unix ms timestamp).
	LastLogin int64 `json:"lastLogin,omitempty"`
}