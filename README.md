# bsctl - A CLI to build collaborative feeds for Bluesky

The purpose of bsctl is to foster micro-communities on Bluesky by making it easy 
for communities to collaboratively create and manage feeds that
discover and rank posts relevant to their communities.

bsctl makes it easy to define feeds in GitHub so that communities can
leverage GitHub for collaboration.

## The Feed Algorithm

The current feed algorithm is pretty simple.

1. The feed retrieve the posts for all handles in a list
   *  [PlatformEngineering List](https://bsky.app/profile/did:plc:5lwweotr4gfb7bbz2fqwdthf/lists/3l7yx65zcse25) is the list for  the [PlatformEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/platformeng)
   * [AIEngineering List](https://bsky.app/profile/did:plc:5lwweotr4gfb7bbz2fqwdthf/lists/3l7z42fommh2l) is the list for the [AIEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/aieng)
   * The source of truth for these lists is in GitHub
1. We filter the posts down to only show posts that include a keyword
   * These keywords are stored in GitHub

## How to Collaborate On Feeds

To collaborate on feeds we store the source of truth for both the list of handles and the keywords in GitHub.
  * The source of truth for the [AIEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/aieng) and
    [PlatformEngineering Feed](https://bsky.app/profile/jeremy.lewi.us/feed/platformeng) feeds is in https://github.com/jlewi/bskylists/

bsctl provides a CLI to make it easy to 
  1. Sync the list of handles from GitHub to a bluesky list
  1. Dump the keywords so they can easily be copied into [Bluesky Feed Creator](https://blueskyfeedcreator.com/)
    * [Bluesky Feed Creator](https://blueskyfeedcreator.com/) is what we currently use to create and serve the feeds


There are two ways to collaborate and improve the algorithm

1. Contribute a PR to update the list of handles or keywords
1. Contribute code to bsctl to improve the algorithm programmatically.
   * For an example checkout [jlewi/bsctl#6](https://github.com/jlewi/bsctl/issues/6)

## How To Setup A New Feed

1. Create a GitHub Repository and create two YAML files 
   1. To define the list of handles
   1. To define the keywords
   * For an example checkout [jlewi/bskylists](https://github.com/jlewi/bskylists/)
1. Use [Blusesky feed creator](https://blueskyfeedcreator.com/) to create the feed
1. Use bsctl to sync the list of handles and keywords to Bluesky
   * Note [jlewi/bsctl#7](https://github.com/jlewi/bsctl/issues/7) is a blocker to creating a new feed
1. Use bsctl to dump the keywords to stdout and then copy and past them into [Bluesky Feed Creator](https://blueskyfeedcreator.com/)   


## Commands For AIEngineering and PlatformEngineering

```bash {"id":"01JBTAGCHJ9H72SM46CHAPGPY3"}
cd ~/git_bskylists
git fetch origin
git checkout origin/main
cd ~/git_bsky
make build
.build/bsctl apply ~/git_bskylists/aiengineering.yaml
.build/bsctl apply ~/git_bskylists/platformengineering.yaml
```

## License

MIT

## Acknowledgement

Originally based on [mattn/bsky](https://github.com/mattn/bsky).
The purpose of this fork is to make the CLI runnable as a client side
web application using [WebAssembly](https://webassembly.org/) and
[goapp](https://github.com/maxence-charriere/go-app).
