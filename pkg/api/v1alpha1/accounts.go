package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	AccountListKind = "AccountList"
	AccountListGVK  = schema.FromAPIVersionAndKind(Group+"/"+Version, AccountListKind)
)

// AccountList is a data structure to hold a list of folks to follow
type AccountList struct {
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`

	// DID is the Decentralized Identifier for the list
	// TOOD(jeremy):
	DID string `json:"did" yaml:"did"`
	// TODO(jeremy): Get rid of this; use Accounts
	Accounts []Account `json:"accounts" yaml:"accounts"`

	// TODO(jeremy): Is it better to have include and exclude lists or just a single list with a field to indicate
	// whether to include or exclude?

	Members []Membership `json:"members" yaml:"members"`
	Exclude []Membership `json:"exclude" yaml:"exclude"`
}

type Membership struct {
	Account Account `json:"account" yaml:"account"`
	// Explanation is a string explaining why the account is in the list
	Explanation string `json:"reason" yaml:"reason"`
}

type Account struct {
	Handle string `json:"handle" yaml:"handle"`
	// DID is the Decentralized Identifier for the account
	// If the DID is specified it will be used and Handle will be ignored as Handles are mutable but DIDs are not.
	DID string `json:"did" yaml:"did"`
}
