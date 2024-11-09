package lists

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jlewi/bsctl/pkg/api/v1alpha1"
	"os"
	"path/filepath"
	"testing"
)

func Test_Prompt(t *testing.T) {
	// Perhaps of this unittest is to verify that the template gets rendered the way we expect
	type testCase struct {
		args         PromptInput
		expectedFile string
	}

	cases := []testCase{
		// No examples
		{
			args: PromptInput{

				Definition: v1alpha1.CommunityDefinition{
					Name: "Platform Engineer",
					Criterion: []string{
						"They are working on an internal developer platform",
						"They describe their job role as platform engineer, ml platform engineer, devops, infrastructure engineer or SRE",
						"They work with technologies used to build platforms; eg. kubernetes, cloud, argo",
						"They describe practices central to platform engineering; e.g. IAC, configuration, containers, gitops, cicd",
					},
				},
				Profile: "I'm building an ml platform at beta.co",
			},
			expectedFile: "no_examples.txt",
		},
		{
			args: PromptInput{
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
				Profile: "I'm building an ml platform at beta.co",
			},
			expectedFile: "examples.txt",
		},
	}

	updateExpected := (os.Getenv("UPDATE_EXPECTED") != "")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	testDir := filepath.Join(cwd, "test_data")
	for _, c := range cases {
		t.Run(c.expectedFile, func(t *testing.T) {

			actual, err := buildPrompt(c.args.Definition, c.args.Profile)
			if err != nil {
				t.Fatalf("Failed to build prompt: %v", err)
			}
			expectedFile := filepath.Join(testDir, c.expectedFile)

			if updateExpected {
				t.Logf("Updating expected file %v", expectedFile)
				if err := os.WriteFile(expectedFile, []byte(actual), 0644); err != nil {
					t.Fatalf("Failed to write expected file: %v", err)
				}
			}

			expected, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("Failed to read expected file: %v", err)
			}

			if d := cmp.Diff(string(expected), actual); d != "" {
				t.Errorf("Unexpected diff:\n%s", d)
			}
		})
	}
}
