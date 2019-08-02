/*
 * Woleet API
 *
 * Welcome to **Woleet API reference documentation**.<br> It is highly recommanded to read the chapters **[introducing Woleet API concepts](https://doc.woleet.io/reference)** before reading this documentation. 
 *
 * API version: 1.6.2
 * Contact: contact@woleet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package woleetapi

type Users struct {
	// Array of users matching the search criteria.
	Content []User `json:"content,omitempty"`
	// `true` if this is the first page. 
	First *bool `json:"first,omitempty"`
	// `true` if this is the last page. 
	Last *bool `json:"last,omitempty"`
	// Number of users in the retrieved page.
	NumberOfElements int32 `json:"numberOfElements,omitempty"`
	// Number of users per page.
	Size int32 `json:"size,omitempty"`
	// Index of the retrieved page (from 0).
	Number int32 `json:"number,omitempty"`
}
