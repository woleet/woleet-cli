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
// KeyStatusEnum The status of the key:<br> - `active`: the key is active: it can be used to sign<br> - `blocked`: the key is blocked: it cannot be used to sign<br> - `revoked` the key is revoked: it will no longer be used to sign 
type KeyStatusEnum string

// List of KeyStatusEnum
const (
	KEYSTATUSENUM_ACTIVE KeyStatusEnum = "active"
	KEYSTATUSENUM_BLOCKED KeyStatusEnum = "blocked"
	KEYSTATUSENUM_REVOKED KeyStatusEnum = "revoked"
)
