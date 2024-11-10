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

	// TODO(jeremy): Is it better to have include and exclude lists or just a single list with a field to indicate
	// whether to include or exclude? I think it might be better to have it as a single list. As the algorithm
	// evolves we might want to reprocess accounts to see whether they are members or not. So we might want to
	// flip the membership status back and forth. This will probably be easier then moving them between and
	// include and exclude list.
	Items []Membership `json:"items" yaml:"items"`

	// TODO(jeremy): Get rid of this; use Accounts
	Accounts []Account `json:"accounts,omitempty" yaml:"accounts,omitempty"`

	// TODO(jeremy): Get rid of these fields
	Members []Membership `json:"members,omitempty" yaml:"members,omitempty"`
	Exclude []Membership `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

type Membership struct {
	Account Account `json:"account" yaml:"account"`
	// Explanation is a string explaining why the account is in the list
	Explanation string `json:"reason,omitempty" yaml:"reason,omitempty"`

	// Whether to include or exclude the account
	Member bool `json:"member" yaml:"member"`
}

type Account struct {
	Handle string `json:"handle,omitempty" yaml:"handle,omitempty"`
	// DID is the Decentralized Identifier for the account
	// If the DID is specified it will be used and Handle will be ignored as Handles are mutable but DIDs are not.
	DID string `json:"did,omitempty" yaml:"did,omitempty"`
}
