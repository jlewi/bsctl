# bsctl

bluesky web CLI client written in Go.

bsctl started as a way to enable people to collaborate on starter packs.
bsctl makes it easy to follow all the accounts listed in a YAML file in GitHub such as
[https://github.com/jlewi/bskylists/blob/main/platformengineering.yaml](https://github.com/jlewi/bskylists/blob/main/platformengineering.yaml).

# Sync The List

```bash {"id":"01JBTAGCHJ9H72SM46CHAPGPY3"}
cd ~/git_bskylists
git fetch origin
git checkout origin/main
cd ~/git_bsky
make build
.build/bsctl apply ~/git_bskylists/aiengineering.yaml
.build/bsctl apply ~/git_bskylists/platformengineering.yaml
```

## Usage

1. Open the webapp at [https://storage.googleapis.com/bsctl/index.html](https://storage.googleapis.com/bsctl/index.html)

2. Set your handle by entering the following command

```sh
config set handle=<YOUR bluesky handle e.g. alice.bsky.social> 
```

* Click the "enter" button
* Alternatively, you can press the enter key on your keyboard **twice**
   * Having to press enter twice is known issue [bsctl/1](https://github.com/jlewi/bsctl/issues/1)

3. Set your bluesky password by entering the following command

```sh
config set password=<YOUR bluesky password> 
```

* Click the "enter" button

4. You can now run commands for example enter the following command to see your followers

```text
followers
```

5. To follow all the accounts listed in a YAML file such as [platformengineering.yaml](https://github.com/jlewi/bskylists/blob/main/platformengineering.yaml)

```sh
follow https://raw.githubusercontent.com/jlewi/bskylists/main/platformengineering.yaml
```

* **Important** The link should point at the raw version of the file

## Why Not Just Use Bluesky Starter Packs

As far as I can tell it doesn't seem possible to programmatically update starterpacks
[thread](https://bsky.app/profile/eribeiro.bsky.social/post/3l7t6gnvyck2a).

For the tech community, using Git/GitHub as a means to collaborate on starter packs seems like a natural fit.
That requires moving the source of truth into git which requires a way to programmatically interact with Bluesky.
While we can't programmatically update starter packs we can achieve a similar effect by making it easy to
subscribe to a bunch of accounts in a list.

## License

MIT

## Acknowledgement

Originally based on [mattn/bsky](https://github.com/mattn/bsky).
The purpose of this fork is to make the CLI runnable as a client side
web application using [WebAssembly](https://webassembly.org/) and
[goapp](https://github.com/maxence-charriere/go-app).
