package credentials

type VendorCredentials struct {
	Id           string `json:"id,omitempty"`
	Account      string `json:"account,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}
