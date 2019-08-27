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
// UserStatusEnum : User status (a `blocked` user cannot sign).
type UserStatusEnum string

// List of UserStatusEnum
const (
	USERSTATUSENUM_ACTIVE UserStatusEnum = "active"
	USERSTATUSENUM_BLOCKED UserStatusEnum = "blocked"
)