package lists

import (
	"context"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/jlewi/bsctl/pkg/api/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"sort"
)

type AccountListController struct {
	client *xrpc.Client
}

func NewAccountListController(client *xrpc.Client) (*AccountListController, error) {
	return &AccountListController{
		client: client,
	}, nil
}

func (c *AccountListController) TidyNode(ctx context.Context, n *yaml.RNode) (*yaml.RNode, error) {
	f := &v1alpha1.AccountList{}
	if err := n.YNode().Decode(f); err != nil {
		return nil, errors.Wrapf(err, "Failed to decode AccountList")
	}

	if err := c.Tidy(ctx, f); err != nil {
		return nil, err
	}

	if err := n.YNode().Encode(f); err != nil {
		return nil, errors.Wrapf(err, "Failed to encode Feed")
	}
	return n, nil
}

func (c *AccountListController) Tidy(ctx context.Context, l *v1alpha1.AccountList) error {
	// Forward convert the lists as we dedupe the accounts
	accounts := map[string]v1alpha1.Membership{}

	for _, a := range l.Accounts {
		m := v1alpha1.Membership{
			Account: a,
			Member:  true,
		}

		accounts[a.Handle] = m
	}

	for _, m := range l.Members {
		m.Member = true
		accounts[m.Account.Handle] = m
	}

	for _, m := range l.Exclude {
		m.Member = true
		m.Member = false
		accounts[m.Account.Handle] = m
	}

	for _, m := range l.Items {
		accounts[m.Account.Handle] = m
	}

	// Sort the accounts by handle
	keys := make([]string, 0, len(accounts))
	for k := range accounts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	members := make([]v1alpha1.Membership, 0, len(keys))
	for _, k := range keys {
		members = append(members, accounts[k])
	}

	l.Items = members

	l.Accounts = nil
	l.Exclude = nil
	l.Members = nil
	return nil
}

func (c *AccountListController) ReconcileNode(ctx context.Context, n *yaml.RNode) error {
	list := &v1alpha1.AccountList{}
	if err := n.YNode().Decode(list); err != nil {
		return errors.Wrapf(err, "Failed to decode AccountList")
	}

	return c.Reconcile(ctx, list)
}

func (c *AccountListController) Reconcile(ctx context.Context, list *v1alpha1.AccountList) error {
	if list.DID == "" {
		return errors.New("List did must be specified. We currently don't support creating new lists yet.")
	}
	if err := AddAllToList(c.client, list.DID, *list); err != nil {
		return err
	}
	return nil
}
