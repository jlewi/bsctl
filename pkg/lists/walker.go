package lists

import (
	"context"
	"encoding/json"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/jlewi/bsctl/pkg/api/v1alpha1"
	"github.com/jlewi/bsctl/pkg/util"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"gopkg.in/yaml.v3"
	"os"
	kYaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

type GraphWalker struct {
	xClient *xrpc.Client
	oClient *openai.Client
}

func (g *GraphWalker) ReconcileNode(ctx context.Context, n *kYaml.RNode) error {
	spec := &v1alpha1.CommunityBuilder{}
	if err := n.YNode().Decode(spec); err != nil {
		return errors.Wrapf(err, "Failed to decode CommunityBuilder")
	}

	return g.Reconcile(ctx, spec)
}

func (g *GraphWalker) TidyNode(ctx context.Context, n *kYaml.RNode) (*kYaml.RNode, error) {
	return n, nil
}

type ClassifyOutput struct {
	Member      bool   `json:"member"`
	Explanation string `json:"explanation"`
}

func NewWalker(xClient *xrpc.Client, oClient *openai.Client) (*GraphWalker, error) {
	return &GraphWalker{
		xClient: xClient,
		oClient: oClient,
	}, nil
}

func (g *GraphWalker) Reconcile(ctx context.Context, buildSpec *v1alpha1.CommunityBuilder) error {
	log := util.LogFromContext(ctx)

	// Is a list of handles we've already seen and made a decision about whether to include them
	// in the community so we don't need to reprocess them
	seen := make(map[string]bool)

	if buildSpec.OutputFile == "" {
		return errors.New("OutputFile must be set")
	}

	// Read the existing account list if one exists.
	accountList := &v1alpha1.AccountList{}

	contents, err := os.ReadFile(buildSpec.OutputFile)

	if err == nil {
		if err := yaml.Unmarshal(contents, accountList); err != nil {
			return errors.Wrapf(err, "Failed to unmarshal account list from file %s", buildSpec.OutputFile)
		}
	} else {
		if !os.IsNotExist(err) {
			return errors.Wrapf(err, "Failed to read file %s", buildSpec.OutputFile)
		}
	}

	// Add the list of seen accounts to the account list
	for _, member := range accountList.Members {
		seen[member.Account.Handle] = true
	}

	for _, member := range accountList.Exclude {
		seen[member.Account.Handle] = true
	}

	for _, seed := range buildSpec.Seeds {
		// Seeds are already in the community so we should mark them as true
		seen[seed.Handle] = true
		if err := g.getFollowers(ctx, buildSpec.Definition, seed.Handle, seen, accountList, buildSpec.OutputFile); err != nil {
			log.Error(err, "Failed to get followers")
			return errors.Wrapf(err, "Failed to get followers for seed %s", seed.Handle)
		}
	}
	return nil
}

func (g *GraphWalker) getFollowers(ctx context.Context, definition v1alpha1.CommunityDefinition, handle string, seen map[string]bool, accountList *v1alpha1.AccountList, outputFile string) error {
	var cursor string
	log := util.LogFromContext(ctx)
	for {

		// Limit controls how often we persist the results to a file since we persist the results after processing
		// each batch of followers
		limit := int64(50)
		followers, err := bsky.GraphGetFollowers(context.TODO(), g.xClient, handle, cursor, limit)
		if err != nil {
			return errors.Wrapf(err, "getting followers for handle %s", handle)
		}

		for _, f := range followers.Followers {
			// If we've already seen this handle then we don't need to process it
			if _, ok := seen[f.Handle]; ok {
				continue
			}

			if f.Description == nil {
				member := v1alpha1.Membership{
					Account: v1alpha1.Account{
						Handle: f.Handle,
						DID:    f.Did,
					},
					Explanation: "The account has no profile description.",
				}
				accountList.Exclude = append(accountList.Exclude, member)
				continue
			}
			prompt, err := buildPrompt(definition, *f.Description)
			if err != nil {
				return errors.Wrapf(err, "Failed to build prompt")
			}
			req := &openai.ChatCompletionRequest{
				Model: openai.GPT4oMini,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				ResponseFormat: &openai.ChatCompletionResponseFormat{
					Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
					JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
						Name: "ClassifyOutput",
						Schema: &jsonschema.Definition{
							Type:        jsonschema.Object,
							Description: "",
							Enum:        nil,
							// TODO(
							Properties: map[string]jsonschema.Definition{
								"member": {
									Type: jsonschema.Boolean,
								},
								"explanation": {
									Type: jsonschema.String,
								},
							},
							Required: []string{"member", "explanation"},
						},
					},
				},
			}

			resp, err := g.oClient.CreateChatCompletion(ctx, *req)
			if err != nil {
				log.Error(err, "OpenAI request failed")
				return errors.Wrapf(err, "Failed to create chat completion")
			}

			if len(resp.Choices) == 0 {
				return errors.New("No choices in response")
			}

			choice := resp.Choices[0]
			output := &ClassifyOutput{}
			if err := json.Unmarshal([]byte(choice.Message.Content), output); err != nil {
				// TODO(jeremy): Should we just keep going? Will this happen for all users
				return errors.Wrap(err, "Failed to unmarshal output")
			}

			member := v1alpha1.Membership{
				Account: v1alpha1.Account{
					Handle: f.Handle,
					DID:    f.Did,
				},
				Explanation: output.Explanation,
			}
			if output.Member {
				accountList.Members = append(accountList.Members, member)
			} else {
				accountList.Exclude = append(accountList.Exclude, member)
			}
		}

		if err := g.save(ctx, accountList, outputFile); err != nil {
			return errors.Wrapf(err, "Failed to save account list")
		}
		if followers.Cursor == nil {
			break
		}
		cursor = *followers.Cursor
	}
	return nil
}

func (g *GraphWalker) save(ctx context.Context, accountList *v1alpha1.AccountList, outFile string) error {
	contents, err := yaml.Marshal(accountList)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal account list")
	}

	if err := os.WriteFile(outFile, contents, 0644); err != nil {
		return errors.Wrapf(err, "Failed to write file %s", outFile)
	}
	return nil
}
