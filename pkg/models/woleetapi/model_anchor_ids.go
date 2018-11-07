/*
 * Woleet API
 *
 * Welcome to **Woleet API reference documentation**.<br> It is highly recommanded to read the chapters **[introducing Woleet API concepts](https://doc.woleet.io/v1.5.1/reference)** before reading this documentation. 
 *
 * API version: 1.5.1
 * Contact: contact@woleet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package woleetapi

type AnchorIds struct {
	// Array of identifiers of anchors matching the search criteria.
	Content []string `json:"content,omitempty"`
	// `true` if this is the first page. 
	First *bool `json:"first,omitempty"`
	// `true` if this is the last page. 
	Last *bool `json:"last,omitempty"`
	// Total number of pages available.
	TotalPages int32 `json:"totalPages,omitempty"`
	// Total number of anchors matching the search criteria.
	TotalElements int32 `json:"totalElements,omitempty"`
	// Number of anchors in the retrieved page.
	NumberOfElements int32 `json:"numberOfElements,omitempty"`
	// Number of anchors per page.
	Size int32 `json:"size,omitempty"`
	// Index of the retrieved page (from 0).
	Number int32 `json:"number,omitempty"`
}