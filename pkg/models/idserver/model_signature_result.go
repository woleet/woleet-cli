/*
Woleet.ID Server

This is Woleet.ID Server API documentation.

API version: 1.2.8
Contact: contact@woleet.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package idserver

import (
	"encoding/json"
)

// SignatureResult struct for SignatureResult
type SignatureResult struct {
	// Public key used to sign (must be the same as the `pubKey` parameter if provided).
	PubKey *string `json:"pubKey,omitempty"`
	// SHA256 hash that is signed (same as the `hashToSign` parameter).
	SignedHash *string `json:"signedHash,omitempty"`
	// Message that is signed (same as the `messageToSign` parameter).
	SignedMessage *string `json:"signedMessage,omitempty"`
	// Signature of `signedMessage` or `signedHash` using the public key `pubKey`, or signature of SHA256(`signedMessage` or `signedHash` + `signedIdentity` + `signedIssuerDomain`) if the identity of the signer and the domain of the identity issuer are included to the signed data. 
	Signature *string `json:"signature,omitempty"`
	// Public URL of the **Identity endpoint** (ie. the URL that anyone can use to get the identity associated to a public key). 
	IdentityURL *string `json:"identityURL,omitempty"`
	// Identity of the signer (as a X500 Distinguished Name).<br> Returned only if `identityToSign` is used. 
	SignedIdentity *string `json:"signedIdentity,omitempty"`
	// Domain of the identity issuer (ie. of the organization who verified the identity).<br> Returned only if `identityToSign` is used. 
	SignedIssuerDomain *string `json:"signedIssuerDomain,omitempty"`
}

// NewSignatureResult instantiates a new SignatureResult object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSignatureResult() *SignatureResult {
	this := SignatureResult{}
	return &this
}

// NewSignatureResultWithDefaults instantiates a new SignatureResult object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSignatureResultWithDefaults() *SignatureResult {
	this := SignatureResult{}
	return &this
}

// GetPubKey returns the PubKey field value if set, zero value otherwise.
func (o *SignatureResult) GetPubKey() string {
	if o == nil || o.PubKey == nil {
		var ret string
		return ret
	}
	return *o.PubKey
}

// GetPubKeyOk returns a tuple with the PubKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SignatureResult) GetPubKeyOk() (*string, bool) {
	if o == nil || o.PubKey == nil {
		return nil, false
	}
	return o.PubKey, true
}

// HasPubKey returns a boolean if a field has been set.
func (o *SignatureResult) HasPubKey() bool {
	if o != nil && o.PubKey != nil {
		return true
	}

	return false
}

// SetPubKey gets a reference to the given string and assigns it to the PubKey field.
func (o *SignatureResult) SetPubKey(v string) {
	o.PubKey = &v
}

// GetSignedHash returns the SignedHash field value if set, zero value otherwise.
func (o *SignatureResult) GetSignedHash() string {
	if o == nil || o.SignedHash == nil {
		var ret string
		return ret
	}
	return *o.SignedHash
}

// GetSignedHashOk returns a tuple with the SignedHash field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SignatureResult) GetSignedHashOk() (*string, bool) {
	if o == nil || o.SignedHash == nil {
		return nil, false
	}
	return o.SignedHash, true
}

// HasSignedHash returns a boolean if a field has been set.
func (o *SignatureResult) HasSignedHash() bool {
	if o != nil && o.SignedHash != nil {
		return true
	}

	return false
}

// SetSignedHash gets a reference to the given string and assigns it to the SignedHash field.
func (o *SignatureResult) SetSignedHash(v string) {
	o.SignedHash = &v
}

// GetSignedMessage returns the SignedMessage field value if set, zero value otherwise.
func (o *SignatureResult) GetSignedMessage() string {
	if o == nil || o.SignedMessage == nil {
		var ret string
		return ret
	}
	return *o.SignedMessage
}

// GetSignedMessageOk returns a tuple with the SignedMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SignatureResult) GetSignedMessageOk() (*string, bool) {
	if o == nil || o.SignedMessage == nil {
		return nil, false
	}
	return o.SignedMessage, true
}

// HasSignedMessage returns a boolean if a field has been set.
func (o *SignatureResult) HasSignedMessage() bool {
	if o != nil && o.SignedMessage != nil {
		return true
	}

	return false
}

// SetSignedMessage gets a reference to the given string and assigns it to the SignedMessage field.
func (o *SignatureResult) SetSignedMessage(v string) {
	o.SignedMessage = &v
}

// GetSignature returns the Signature field value if set, zero value otherwise.
func (o *SignatureResult) GetSignature() string {
	if o == nil || o.Signature == nil {
		var ret string
		return ret
	}
	return *o.Signature
}

// GetSignatureOk returns a tuple with the Signature field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SignatureResult) GetSignatureOk() (*string, bool) {
	if o == nil || o.Signature == nil {
		return nil, false
	}
	return o.Signature, true
}

// HasSignature returns a boolean if a field has been set.
func (o *SignatureResult) HasSignature() bool {
	if o != nil && o.Signature != nil {
		return true
	}

	return false
}

// SetSignature gets a reference to the given string and assigns it to the Signature field.
func (o *SignatureResult) SetSignature(v string) {
	o.Signature = &v
}

// GetIdentityURL returns the IdentityURL field value if set, zero value otherwise.
func (o *SignatureResult) GetIdentityURL() string {
	if o == nil || o.IdentityURL == nil {
		var ret string
		return ret
	}
	return *o.IdentityURL
}

// GetIdentityURLOk returns a tuple with the IdentityURL field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SignatureResult) GetIdentityURLOk() (*string, bool) {
	if o == nil || o.IdentityURL == nil {
		return nil, false
	}
	return o.IdentityURL, true
}

// HasIdentityURL returns a boolean if a field has been set.
func (o *SignatureResult) HasIdentityURL() bool {
	if o != nil && o.IdentityURL != nil {
		return true
	}

	return false
}

// SetIdentityURL gets a reference to the given string and assigns it to the IdentityURL field.
func (o *SignatureResult) SetIdentityURL(v string) {
	o.IdentityURL = &v
}

// GetSignedIdentity returns the SignedIdentity field value if set, zero value otherwise.
func (o *SignatureResult) GetSignedIdentity() string {
	if o == nil || o.SignedIdentity == nil {
		var ret string
		return ret
	}
	return *o.SignedIdentity
}

// GetSignedIdentityOk returns a tuple with the SignedIdentity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SignatureResult) GetSignedIdentityOk() (*string, bool) {
	if o == nil || o.SignedIdentity == nil {
		return nil, false
	}
	return o.SignedIdentity, true
}

// HasSignedIdentity returns a boolean if a field has been set.
func (o *SignatureResult) HasSignedIdentity() bool {
	if o != nil && o.SignedIdentity != nil {
		return true
	}

	return false
}

// SetSignedIdentity gets a reference to the given string and assigns it to the SignedIdentity field.
func (o *SignatureResult) SetSignedIdentity(v string) {
	o.SignedIdentity = &v
}

// GetSignedIssuerDomain returns the SignedIssuerDomain field value if set, zero value otherwise.
func (o *SignatureResult) GetSignedIssuerDomain() string {
	if o == nil || o.SignedIssuerDomain == nil {
		var ret string
		return ret
	}
	return *o.SignedIssuerDomain
}

// GetSignedIssuerDomainOk returns a tuple with the SignedIssuerDomain field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SignatureResult) GetSignedIssuerDomainOk() (*string, bool) {
	if o == nil || o.SignedIssuerDomain == nil {
		return nil, false
	}
	return o.SignedIssuerDomain, true
}

// HasSignedIssuerDomain returns a boolean if a field has been set.
func (o *SignatureResult) HasSignedIssuerDomain() bool {
	if o != nil && o.SignedIssuerDomain != nil {
		return true
	}

	return false
}

// SetSignedIssuerDomain gets a reference to the given string and assigns it to the SignedIssuerDomain field.
func (o *SignatureResult) SetSignedIssuerDomain(v string) {
	o.SignedIssuerDomain = &v
}

func (o SignatureResult) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.PubKey != nil {
		toSerialize["pubKey"] = o.PubKey
	}
	if o.SignedHash != nil {
		toSerialize["signedHash"] = o.SignedHash
	}
	if o.SignedMessage != nil {
		toSerialize["signedMessage"] = o.SignedMessage
	}
	if o.Signature != nil {
		toSerialize["signature"] = o.Signature
	}
	if o.IdentityURL != nil {
		toSerialize["identityURL"] = o.IdentityURL
	}
	if o.SignedIdentity != nil {
		toSerialize["signedIdentity"] = o.SignedIdentity
	}
	if o.SignedIssuerDomain != nil {
		toSerialize["signedIssuerDomain"] = o.SignedIssuerDomain
	}
	return json.Marshal(toSerialize)
}

type NullableSignatureResult struct {
	value *SignatureResult
	isSet bool
}

func (v NullableSignatureResult) Get() *SignatureResult {
	return v.value
}

func (v *NullableSignatureResult) Set(val *SignatureResult) {
	v.value = val
	v.isSet = true
}

func (v NullableSignatureResult) IsSet() bool {
	return v.isSet
}

func (v *NullableSignatureResult) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSignatureResult(val *SignatureResult) *NullableSignatureResult {
	return &NullableSignatureResult{value: val, isSet: true}
}

func (v NullableSignatureResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSignatureResult) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


