# How to release a new version of Trento

> Note: this document is a draft!

## Pre-requisites

Install [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator) globally in your dev box:
```
gem install github_changelog_generator
```

## Update the changelog and create a new tag

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
