package lists

import (
	_ "embed"
	"github.com/jlewi/bsctl/pkg/api/v1alpha1"
	"github.com/pkg/errors"
	"strings"
	"text/template"
)

//go:embed profile_prompt.tmpl
var promptTemplateString string

var (
	promptTemplate = template.Must(template.New("prompt").Parse(promptTemplateString))
)

type PromptInput struct {
	Definition v1alpha1.CommunityDefinition
	Profile    string
}

func buildPrompt(definition v1alpha1.CommunityDefinition, profile string) (string, error) {
	var sb strings.Builder
	input := PromptInput{
		Definition: definition,
		Profile:    profile,
	}
	if err := promptTemplate.Execute(&sb, input); err != nil {
		return "", errors.Wrapf(err, "Failed to execute prompt template")
	}
	return sb.String(), nil
}
