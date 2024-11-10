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