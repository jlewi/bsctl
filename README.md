# bsctl - A CLI to build collaborative feeds for Bluesky

The purpose of bsctl is to foster micro-communities on Bluesky by making it easy
for communities to collaboratively create and manage feeds that
discover and rank posts relevant to their communities.

bsctl makes it easy to define feeds in GitHub so that communities can
leverage GitHub for collaboration.

## The Feed Algorithm

The current feed algorithm is pretty simple.

1. The feed retrieve the posts for all handles in a list

   * [PlatformEngineering List](https://bsky.app/profile/did:plc:5lwweotr4gfb7bbz2fqwdthf/lists/3l7yx65zcse25) is the list for  the [PlatformEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/platformeng)
   * [AIEngineering List](https://bsky.app/profile/did:plc:5lwweotr4gfb7bbz2fqwdthf/lists/3l7z42fommh2l) is the list for the [AIEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/aieng)
   * The source of truth for these lists is in GitHub

2. We filter the posts down to only show posts that include a keyword

   * These keywords are stored in GitHub

## How to Collaborate On Feeds

To collaborate on feeds we store the source of truth for both the list of handles and the keywords in GitHub.

* The source of truth for the [AIEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/aieng) and
   [PlatformEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/platformeng) feeds is in https://github.com/jlewi/bskylists/

bsctl provides a CLI to make it easy to

1. Sync the list of handles from GitHub to a bluesky list
2. Dump the keywords so they can easily be copied into [Bluesky Feed Creator](https://blueskyfeedcreator.com/)
   * [Bluesky Feed Creator](https://blueskyfeedcreator.com/) is what we currently use to create and serve the feeds

There are two ways to collaborate and improve the algorithm

1. Contribute a PR to update the list of handles or keywords
2. Contribute code to bsctl to improve the algorithm programmatically.
   * For an example checkout [jlewi/bsctl#6](https://github.com/jlewi/bsctl/issues/6)

## How To Setup A New Feed

1. Create a GitHub Repository and create two YAML files

   1. To define the list of handles
   2. To define the keywords

   * For an example checkout [jlewi/bskylists](https://github.com/jlewi/bskylists/)

2. Use [Blusesky feed creator](https://blueskyfeedcreator.com/) to create the feed
3. Use bsctl to sync the list of handles and keywords to Bluesky

   * Note [jlewi/bsctl#7](https://github.com/jlewi/bsctl/issues/7) is a blocker to creating a new feed

4. Use bsctl to dump the keywords to stdout and then copy and past them into [Bluesky Feed Creator](https://blueskyfeedcreator.com/)

## Identifying Community Members Algorithmically

Rather than manually curating a list of handles, we can use algorithms to identify community members.
The algorithm is very simple

1. Take a set of seed handles
1. Fetch the profile for everyone who **follows** those accounts
1. Use an LLM to decide whether the profile satisfies some criterion

To use this functionality define a `CommunityBuilder` resource in a YAML file like the one below.

```yaml
apiVersion: "bsctl.dev/v1alpha1"
kind: "CommunityBuilder"
metadata:
    labels: {}
    name: PlatformCommunity
definition:
    name: Platform Engineer
    criterion:
        - They are working on an internal developer platform
        - They describe their job role as platform engineer, ml platform engineer, devops, infrastructure engineer or SRE
        - They work with technologies used to build platforms; eg. kubernetes, cloud, argo
        - They describe practices central to platform engineering; e.g. IAC, configuration, containers, gitops, cicd
    examples:
        - profile: I'm a platform engineer at acme.co
          member: true
          explanation: ""
seeds:
    - handle: kelseyhightower.com
      did: ""
outputFile: /Users/jlewi/git_bskylists//Users/jlewi/git_bskylists/kelseyhightower.followers.platformengineering.yaml
```

* Everything in the `definition` section is used as input to the [LLM Prompt](https://github.com/jlewi/bsctl/blob/main/pkg/lists/profile_prompt.tmpl)
* To perform the graph walk first configure your OpenAI APIKey
  
  ```
  bsctl config set openai.apiKeyFile=/path/to/openai/apikey
  ```

* You can then perform the graph walk using the CLI

  ```
  bsctl apply /path/to/your/community.yaml
  ```

* This will write the results to the outputFile you specified
* The output will contain an AccountList that tells you
  1. Whether the account is a member of the community
  1. A reason for the decision

* For an example of the output see [kelseyhightower.followers.platformengineering.yaml]

## License

MIT

## Acknowledgement

Originally based on [mattn/bsky](https://github.com/mattn/bsky).
The purpose of this fork is to make the CLI runnable as a client side
web application using [WebAssembly](https://webassembly.org/) and
[goapp](https://github.com/maxence-charriere/go-app).
