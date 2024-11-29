# Platform Engineering Graph Walk

* Lets Update Platform Engineering Using the First verison of the graph walker
* We need to tidy the lists to use the new format

```sh
cd ..
make build
```

```sh
../.build/bsctl tidy /Users/jlewi/git_bskylists/aiengineering.yaml
```

```bash {"id":"01JC9V08MWAKR148RG5MA2R2HH"}
../.build/bsctl tidy /Users/jlewi/git_bskylists/platformengineering.yaml
```

# Merge Kelsey's list into the platform engineering list

* How big is the group currently

```bash
yq ".items | length" /Users/jlewi/git_bskylists/platformengineering.yaml
```

```bash {"id":"01JC9VSV7Z96A62PERW4HKY5KT"}
../.build/bsctl merge --dest=/Users/jlewi/git_bskylists/platformengineering.yaml \
  --source=/Users/jlewi/git_bskylists/kelseyhightower.followers.platformengineering.yaml \
  --source-filter=members
```

```bash
yq ".items | length" /Users/jlewi/git_bskylists/platformengineering.yaml
```

* So it had 204 items before and now it has 548
* Now lets sync it to bluesky

```bash
../.build/bsctl apply /Users/jlewi/git_bskylists/platformengineering.yaml
```

* Changes don't seem to be reflected in the blueskyfeed generator
* Bluesky feed generator says it was last synced 22 hours ago
   * Lets give it time to sync

```bash
yq '[(.items[] | select(.member == true))] | length' ~/git_bskylists/kelseyhightower.followers.platformengineering.yaml
```

## Lets find my followers

```bash {"id":"01JC9XTX3AT7XEVS6AZ752NWYT"}
../.build/bsctl apply /Users/jlewi/git_bskylists/jeremylewi.followers.communitybuilder.yaml
```

# Lets merge in my followers

```bash
yq ".items | length" /Users/jlewi/git_bskylists/platformengineering.yaml
```

```bash {"id":"01JCAAJ0SKBKGC48KQAQ0YXKAF"}
../.build/bsctl merge --dest=/Users/jlewi/git_bskylists/platformengineering.yaml \
  --source=/Users/jlewi/git_bskylists/jeremylewi.followers.platformengineering.yaml \
  --source-filter=members
```

```bash
yq ".items | length" /Users/jlewi/git_bskylists/platformengineering.yaml
```

Looks like that didn't add any members so looks like everyone who's following me that's a platform engineer is already there

* How many platform engineers do I have following me

```bash {"id":"01JCAAMYQ15YRFHYP6BM4NQRG9"}
yq '[(.items[] | select(.member == true))] | length' /Users/jlewi/git_bskylists/jeremylewi.followers.platformengineering.yaml
```

I only have 15 followers that are in the platform engineering community

* Lets get there handles

```bash {"id":"01JCAARBRD24KE09A44MS3A2Q0"}
yq '(.items[] | select(.member == true)) | .account.handle' /Users/jlewi/git_bskylists/jeremylewi.followers.platformengineering.yaml
```

## Lets redo my followers

* I've picked up a bunch in the base couple hours

```bash
../.build/bsctl apply /Users/jlewi/git_bskylists/jeremylewi.followers.communitybuilder.yaml
```

* So now I have 18 as opposed to 15 followers that meet the criterion.

```bash

```