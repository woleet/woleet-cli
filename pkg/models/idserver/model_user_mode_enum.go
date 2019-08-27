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
// UserModeEnum : User mode, esign the user is a real user, and his keys are used to create electronic signatures. seal the user represents an organization, and his keys are used to create server seals. 
type UserModeEnum string

// List of UserModeEnum
const (
	USERMODEENUM_SEAL UserModeEnum = "seal"
	USERMODEENUM_ESIGN UserModeEnum = "esign"
)