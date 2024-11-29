Sync the lists 

```bash
cd ~/git_bskyctl/.build
git clone git@github.com:jlewi/bskylists.git 
```

```bash
# TODO(jeremy): We should checkout a seperate copy of the bskylists repo into .build
cd ../.build/bskylists
git fetch origin
git checkout origin/main
```

```bash {"id":"01JCVXFYBQY2YVJ175REQDYCER"}
cd ..
make build
```

```bash {"id":"01JCVX3AACGCBTQA817GJYXFBP"}
../.build/bsctl apply ../.build/bskylists/aiengineering.yaml
../.build/bsctl apply ../.build/bskylists/platformengineering.yaml
```