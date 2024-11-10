package v1alpha1

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	CommunityBuilderKind = "CommunityBuilder"
	CommunityBuilderGVK  = schema.FromAPIVersionAndKind(Group+"/"+Version, CommunityBuilderKind)
)

// TODO(jeremy): I think the name should be a noun not a verb.

type CommunityBuilder struct {
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`

	// Definition is the definition of the community. It parameterizes the LLM prompt to classifiy accounts
	Definition CommunityDefinition `json:"definition" yaml:"definition"`
	// Seeds is a list of accounts to seed the graph with
	Seeds []Account `json:"seeds" yaml:"seeds"`

	// OutputFile is the file to write the AccountList to
	OutputFile string `json:"outputFile" yaml:"outputFile"`
}

type CommunityDefinition struct {
	// Name is the name of the community
	Name string `json:"name" yaml:"name"`

	// Criterion is a list of criterion for including accounts in the community
	Criterion []string `json:"criterion" yaml:"criterion"`

	// Example is a list of examples to help classify accounts
	Examples []ProfileExample `json:"examples" yaml:"examples"`
}

type ProfileExample struct {
	// Profile is the example profile
	Profile string `json:"profile" yaml:"profile"`
	// Member is true or false depending on whether the profile is a member of the community
	Member bool `json:"member" yaml:"member"`

	// Explanation is the explanation of why the profile is a member or not
	Explanation string `json:"explanation" yaml:"explanation"`
}
