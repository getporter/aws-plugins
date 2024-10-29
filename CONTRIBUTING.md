# Contributing Guide

This is part of the [Porter][porter] project. If you are a new contributor,
check out our [New Contributor Guide][new-contrib]. The Porter [Contributing
Guide][contrib] also has lots of information about how to interact with the
project.

[porter]: https://github.com/getporter/porter
[new-contrib]: https://porter.sh/contribute
[contrib]: https://porter.sh/src/CONTRIBUTING.md

---

* [Initial setup](#initial-setup)
* [Magefile explained](#magefile-explained)

---

# Initial setup

You need to have [porter installed](https://porter.sh/install) first. Then run
`mage build install`. This will build and install the plugin into your porter
home directory.

## Magefile explained

We use [mage](https://magefile.org) instead of make. If you don't have mage installed already,
you can install it with `go run mage.go EnsureMage`.

[mage]: https://magefile.org

Mage targets are not case-sensitive, but in our docs we use camel case to make
it easier to read. You can run either `mage Build` or `mage build` for
example.

Run `mage` without any arguments to see a list of the available targets.
Below are some commonly used targets:

* `Build` builds the plugin.
* `Install` installs the plugin into **~/.porter/plugins**.
* `TestUnit` runs the unit tests.
