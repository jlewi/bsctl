package e2etests

import (
	"context"
	"github.com/jlewi/bsctl/pkg/api/v1alpha1"
	"github.com/jlewi/bsctl/pkg/lists"
	"github.com/jlewi/bsctl/pkg/testutil"
	"os"
	"testing"
)

func Test_Walker(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skipf("Test_AccountsListApply is a manual test that is skipped in CICD")
	}

	stuff, err := testutil.New()
	if err != nil {
		t.Fatalf("testSetup() = %v, wanted nil", err)
	}

	app := stuff.App
	client, err := app.GetOAIClient(context.Background())
	if err != nil {
		t.Fatalf("Failed to create OAI client; error %+v", err)
	}
	w, err := lists.NewWalker(stuff.Client, client)
	if err != nil {
		t.Fatalf("Failed to create walker; %+v", err)
	}

	f, err := os.CreateTemp("", "accounts.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file; %+v", err)
	}

	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close file; %+v", err)
	}

	oFile := f.Name()

	t.Logf("Output file: %s", oFile)

	buildSpec := &v1alpha1.CommunityBuilder{
		OutputFile: oFile,
		Seeds: []v1alpha1.Account{
			{
				Handle: "jeremy.lewi.us",
			},
		},
		Definition: v1alpha1.CommunityDefinition{
			Name: "Platform Engineer",
			Examples: []v1alpha1.ProfileExample{
				{
					Profile: "I'm a platform engineer at acme.co",
					Member:  true,
				},
			},
			Criterion: []string{
				"They are working on an internal developer platform",
				"They describe their job role as platform engineer, ml platform engineer, devops, infrastructure engineer or SRE",
				"They work with technologies used to build platforms; eg. kubernetes, cloud, argo",
				"They describe practices central to platform engineering; e.g. IAC, configuration, containers, gitops, cicd",
			},
		},
	}

	//b, err := yaml.Marshal(buildSpec)
	//if err != nil {
	//	t.Fatalf("Failed to marshal buildSpec; %+v", err)
	//}
	//if err := os.WriteFile("/tmp/platform_community_builder.yaml", b, 0644); err != nil {
	//	t.Fatalf("Failed to write buildSpec; %+v", err)
	//}
	//return

	if err := w.Reconcile(context.Background(), buildSpec); err != nil {
		t.Fatalf("Reconcile() = %v, wanted nil", err)
	}
}
