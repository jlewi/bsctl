package lists

import (
	"context"
	"fmt"
	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/go-logr/zapr"
	"github.com/jlewi/bsctl/pkg"
	"github.com/jlewi/bsctl/pkg/api/v1alpha1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sort"
	"time"
)

const (
	referenceListPurpose = "app.bsky.graph.defs#referencelist"
)

// ListRecord struct to represent the list record structure
type ListRecord struct {
	Type        string    `json:"$type"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Users       []User    `json:"users"`
}

// User struct to represent each user in the list
type User struct {
	DID string `json:"did"`
}

//
//func GetList(client *xrpc.Client, listAtRef string) {
//	//TODO(jeremy): Support the cursor?
//	cursor := ""
//	limit := int64(100)
//	bsky.GraphGetList(context.Background(), client, cursor, limit, listAtRef)
//}

// CreateListRecord sends a request to the PDS server to create a list record.
//
// TODO(jeremy): How should we check if a list of the given name already exists?
func CreateListRecord(client *xrpc.Client, name string, description string) (*comatproto.RepoCreateRecord_Output, error) {
	//var out bytes.Buffer
	//if err := client.Do(context.Background(), xrpc.Procedure, "", "app.bsky.graph.list", nil, record, &out); err != nil {
	//	return err
	//}
	log := zapr.NewLogger(zap.L())

	// I think we need to create a list and then we create GraphListItem
	block := bsky.GraphList{
		LexiconTypeID: "app.bsky.graph.list",
		CreatedAt:     time.Now().Local().Format(time.RFC3339),
		Name:          name,
		Description:   pkg.StringPtr(description),
		Purpose:       pkg.StringPtr(referenceListPurpose),
	}

	//block := bsky.GraphBlock{
	//	LexiconTypeID: "app.bsky.graph.block",
	//	CreatedAt:     time.Now().Local().Format(time.RFC3339),
	//	Subject:       profile.Did,
	//}

	resp, err := comatproto.RepoCreateRecord(context.TODO(), client, &comatproto.RepoCreateRecord_Input{
		Collection: "app.bsky.graph.list",
		Repo:       client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: &block,
		},
	})

	if err != nil {
		log.Error(err, "Failed to create list record")
		return nil, err
	}

	log.Info("List record created", "record", resp)

	return resp, nil
}

//func GetList(client *xrpc.Client, listUri string) (*comatproto.RepoCreateRecord_Output, error) {
//	//var out bytes.Buffer
//	//if err := client.Do(context.Background(), xrpc.Procedure, "", "app.bsky.graph.list", nil, record, &out); err != nil {
//	//	return err
//	//}
//	log := zapr.NewLogger(zap.L())
//
//	//// I think we need to create a list and then we create GraphListItem
//	//block := bsky.GraphList{
//	//	LexiconTypeID: "app.bsky.graph.list",
//	//	CreatedAt:     time.Now().Local().Format(time.RFC3339),
//	//	Name:          name,
//	//	Description:   StringPtr(description),
//	//	Purpose:       StringPtr(referenceListPurpose),
//	//}
//
//	//block := bsky.GraphBlock{
//	//	LexiconTypeID: "app.bsky.graph.block",
//	//	CreatedAt:     time.Now().Local().Format(time.RFC3339),
//	//	Subject:       profile.Did,
//	//}
//
//	bsky.GraphGetList()
//	resp, err := comatproto.G(context.TODO(), client, &comatproto.RepoCreateRecord_Input{
//		Collection: "app.bsky.graph.list",
//		Repo:       client.Auth.Did,
//		Record: &lexutil.LexiconTypeDecoder{
//			Val: &block,
//		},
//	})
//
//	if err != nil {
//		log.Error(err, "Failed to create list record")
//		return nil, err
//	}
//
//	log.Info("List record created", "record", resp)
//
//	return resp, nil
//}

func AddAllToList(client *xrpc.Client, listURI string, source v1alpha1.AccountList) error {
	log := zapr.NewLogger(zap.L())
	for _, m := range source.Items {

		if !m.Member {
			continue
		}

		profile, err := bsky.ActorGetProfile(context.TODO(), client, m.Account.Handle)
		if err != nil {
			var xErr *xrpc.Error
			if errors.As(err, &xErr) {
				if 400 == xErr.StatusCode {
					log.Error(err, "Profile not found for handle", "handle", m.Account.Handle)
					continue
				}
			}
			return fmt.Errorf("cannot get profile: %w", err)
		}
		AddToList(client, listURI, profile.Did)
	}

	return nil

}

// AddToList adds a subjectDid to the list
func AddToList(client *xrpc.Client, listURI string, subjectDid string) error {
	item := bsky.GraphListitem{
		LexiconTypeID: "app.bsky.graph.listitem",
		CreatedAt:     time.Now().Local().Format(time.RFC3339),
		List:          listURI,
		Subject:       subjectDid,
	}

	//block := bsky.GraphBlock{
	//	LexiconTypeID: "app.bsky.graph.block",
	//	CreatedAt:     time.Now().Local().Format(time.RFC3339),
	//	Subject:       profile.Did,
	//}

	itemResp, err := comatproto.RepoCreateRecord(context.TODO(), client, &comatproto.RepoCreateRecord_Input{
		Collection: "app.bsky.graph.listitem",
		Repo:       client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: &item,
		},
	})
	log := zapr.NewLogger(zap.L())
	log.Info("List item record created", "item", itemResp)

	return errors.Wrapf(err, "Failed to add subject to list: %s", listURI)
}

//func MutateList(client *xrpc.Client, listRef string) error {
//	cursor := ""
//	limit := int64(100)
//	list, err := bsky.GraphGetList(context.Background(), client, cursor, limit, listRef)
//
//	if err != nil {
//		return errors.Wrapf(err, "Failed to fetch list: %s", listRef)
//	}
//
//	itemResp, err := comatproto.RepoPutRecord(context.TODO(), client, &comatproto.RepoCreateRecord_Input{
//		Collection: "app.bsky.graph.listitem",
//		Repo:       client.Auth.Did,
//		Record: &lexutil.LexiconTypeDecoder{
//			Val: &item,
//		},
//	})
//}

type ListFilter string

const (
	// IncludeAll includes all members
	IncludeAll ListFilter = "all"
	// IncludeMembers includes items that are members; i.e. Membership.Member is true
	IncludeMembers ListFilter = "members"
	// IncludeNonMembers includes items that are members; i.e. Membership.Member is false
	IncludeNonMembers ListFilter = "nonmembers"
)

// MergeFollowLists computes the union of two lists
func MergeFollowLists(dest *v1alpha1.AccountList, src v1alpha1.AccountList, srcFilter ListFilter) {
	// Use a map to store unique strings from both lists
	uniqueStrings := map[string]v1alpha1.Membership{}

	// Add elements from the second list to the map
	for _, item := range src.Items {
		include := false
		if srcFilter == IncludeAll {
			include = true
		}
		if srcFilter == IncludeMembers && item.Member {
			include = true
		}
		if srcFilter == IncludeNonMembers && !item.Member {
			include = true
		}
		if !include {
			continue
		}
		uniqueStrings[item.Account.Handle] = item
	}

	// Add elements from the dest list to the map
	// By adding dest list after src list, we ensure that the value in dest list overrides the src list value if
	// there are duplicates. I'm not sure if this is the desired behavior.
	for _, item := range dest.Items {
		uniqueStrings[item.Account.Handle] = item
	}

	// Convert map keys to a slice
	result := make([]string, 0, len(uniqueStrings))
	for item := range uniqueStrings {
		result = append(result, item)
	}

	// Sort the result slice
	sort.Strings(result)

	dest.Items = make([]v1alpha1.Membership, 0, len(result))
	for _, item := range result {
		dest.Items = append(dest.Items, uniqueStrings[item])
	}
}

// MergeListToFile merges the users in a list to file.
func MergeListToFile(client *xrpc.Client, sourceFile string, listURI string) error {
	b, err := os.ReadFile(sourceFile)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "cannot read file %s", sourceFile)
	}

	dest := &v1alpha1.AccountList{}
	if err == nil {
		nodes, err := kio.FromBytes(b)
		if err != nil {
			return errors.Wrapf(err, "cannot read file %s", sourceFile)
		}

		node := nodes[0]

		if err := node.YNode().Decode(dest); err != nil {
			return errors.Wrapf(err, "cannot unmarshal AccountList from file %s", sourceFile)
		}
	}
	output, err := DumpList(client, listURI)
	if err != nil {
		return err
	}

	MergeFollowLists(dest, *output, IncludeAll)

	outB, err := yaml.Marshal(dest)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal AccountList to file %s", sourceFile)
	}

	if err := os.WriteFile(sourceFile, outB, 0644); err != nil {
		return errors.Wrapf(err, "cannot write file %s", sourceFile)
	}
	return nil
}

// DumpList dumps the list
func DumpList(client *xrpc.Client, listUri string) (*v1alpha1.AccountList, error) {
	// Get the list
	cursor := ""

	result := &v1alpha1.AccountList{
		Items: make([]v1alpha1.Membership, 0),
	}
	for {
		output, err := bsky.GraphGetList(context.Background(), client, cursor, 100, listUri)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get list associated with the starter pack")
		}

		for _, item := range output.Items {
			result.Items = append(result.Items, v1alpha1.Membership{
				Account: v1alpha1.Account{
					Handle: item.Subject.Handle,
				}})
		}

		if output.Cursor == nil {
			break
		}
		cursor = *output.Cursor
	}
	return result, nil
}
