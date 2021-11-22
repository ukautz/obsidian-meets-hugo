# `omh` - Obsidian Meets Hugo

Command line tool to marry [Obsidian](https://obsidian.md/) vaults to [Hugo](https://gohugo.io/) published websites.

```sh
# convert and copy all obsidian notes in directory (and sub-directories)
#  into hugo path in `/path/to/hugo/content/obsidian` and respective static
#  files into `/path/to/hugo/static/obsidian`
$ omh \
    --obsidian-root /path/to/obsidian \
    --hugo-root /path/to/hugo
```

See `omh -h` for extended options.

_Note: on Mac you can find your iCloud synced notes in `~/Library/Mobile\ Documents/iCloud\~md\~obsidian/Documents/`_

## Install

```sh
$ go install github.com/ukautz/obsidian-meets-hugo/cmds/omh
```

## Use-Case

This command line tool allows you to easily export an Obsidian vault, or a sub-set thereof, into a Hugo published website.

I am using this tool to publish my own notes - that follow a standard in between [Zettelkasten and Wikipedia](https://en.wikipedia.org/wiki/Zettelkasten) - to my [Blog](https://ulrichkautz.com), as you can see here: <https://ulrichkautz.com/zettel/>. This way I have them easily available and can reference them in Blog entries.

## License

MIT

## Alternatives

Things I found that do not exactly fit my needs, but maybe yours:

- <https://github.com/khalednassar/obyde>
- <https://github.com/jackyzha0/hugo-obsidian>
