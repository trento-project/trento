# How to release a new version of Trento

> Note: this document is a draft!

## Pre-requisites

Install [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator) globally in your dev box:
```
gem install github_changelog_generator
```

## Update the changelog and create a new tag

The automatic changelog generation leverages GitHub labels very heavily to produce a meaningful output following the [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) specification, grouping pull requests and issues in sections.

Use the labels as follows:
- `enhancement` or `addition` items go in the `Added` section;
- `bug` or `fix` items go in the `Fixed` section;
- `removal` items go in the `Removed` section;
- unlabelled pull requests go in the `Other Changes` section;
- unlabelled closed issues are ignored.

You don't have to label everything: the intent of the changelog is to communicate highlights to end-users, while also being comprehensive; this is why the `Other changes` section catches all the unlabelled items and is rendered last.

Once you do a quick round of issues/PR triaging to apply labels in a meaningful way, follow these steps:

```bash
# always create a dedicated release branch
git switch -c release-x.y.z

# x1.y1.z1 is the previous release tag
github_changelog_generator --since-tag=x1.y1.z1 --future-release=x.y.z

git add CHANGELOG.md
git commit -m "add x.y.z changelog entry"

# maybe make some other last minute changes
# [...]

# merge and tag, making sure the tag is on the merge commit
git switch main
git merge --no-ff release-x.y.z
git tag x.y.z

# don't forget to force update the rolling tag!
git fetch --tags -f

# push directly
git push --tags origin main
```

Optionally, open a pull request from the release branch instead of tagging and pushing manually.

## GitHub release

> Note: this step will soon be automated.

Go to the [project releases page](https://github.com/trento-project/trento/releases) and create a new release, then:

- use the just created git tag as the release tag and title;
- copy-paste the last changelog entry from `CHANGELOG.md` as the release body;
- hit the green button;
- profit!
