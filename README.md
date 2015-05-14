### gb-vendorfile

`gb-vendorfile` is a plugin for [`gb`](https://github.com/constabulary/gb)
that knows how to read a dependency file similar to python requirements.txt
files.

### Usage

Create a `Vendorfile` like below:

```
# fs
https://github.com/kr/fs/archive/2788f0dbd16903de03cb8186e5c7d97b69ad387b.zip
```

Now we can set up `gb`'s vendor path from this file.

```bash
$ gb plugin vendorfile get
... vendoring github.com/kr.fs
```

You vendor path should now be set up and you can use regular `gb` workflow.

```bash
$ gb build
```

Now be sure to check in the `Vendorfile` and `.gitignore` your `vendor/` folder.

### Why?

I really like how `gb` allows you to have per project GOPATHS without much
setup and configuration. I just really don't like checking in the source of
all of your dependencies. Having a dependency file so we can get the specific
revision will still allow us to have reproducible builds, as long as you have
network connectivity, and your dependency can be downloaded from the internet.

Also, because I can.

### *Note*

This is a POC that only runs on *nix. Use at your own discretion.
