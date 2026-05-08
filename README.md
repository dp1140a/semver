# semver

Yet another semantic versioning tool written in golang that follows the [Semver v2 Spec](https://semver.org/)
I wrote this cause it follows how I like to use these tools: less flags and more subcommands that make sense.

The tool accepts both canonical SemVer strings like `1.0.0-rc.1` and the common `v`-prefixed form like `v1.0.0-rc.1`. The leading `v` is treated as optional compatibility syntax rather than part of the SemVer spec. When mutating an existing `VERSION` file, `semver` preserves the file's current prefix style unless you explicitly provide a full replacement version with a different style.

## Getting Started

### Install
If you already have golang installed you can install by running the command:

```sh
go install github.com/dp1140a/semver@latest
```

### check install
Check if the semver was installed running the command:

```sh
semver 
```

## Usage:
### semver [subcommand]
Run by itself semver will return the current version string. For example if the current version is ```1.2.3```:
```
$ semver --> 1.2.3
```
This is equivalent to running

```bash
$ semver version
$ semver version -f string
```
Usage:
```
semver [flags]
semver [command]
```

Available Commands:
* bump -- Will bump the current version
* completion -- Generate the autocompletion script for the specified shell
* help -- Help about any command
* init -- A brief description of your command
* set -- Set command for PreRelease or Build information
* version -- Prints the current version

Flags:
-h, --help   help for semver. Available to all commands

Use "semver [command] --help" for more information about a command.

---

### init

Will launch an interactive console to launch a semver project.  This must be done in an existing git repo.
It will create a file called VERSION that will be used to track version information.  If an Existing VERSION file is found it will ask if you want to overwrite.

Usage:
```semver init```

---

### bump

If no subcommand is specified this command will bump the Patch version.  For example if our current version is 0.1.0:

```$ semver bump --> 0.1.1```
Is the same as
```$ semver bump patch --> 0.1.1```

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For instance if our current version is 1.2.3-alpha:

```
$ semver bump major --> 2.0.0
$ semver bump minor --> 1.3.0
$ semver bump patch --> 1.2.4
```

Usage:
```
semver bump 
semver bump [command]
```

Available Commands:
major       Will bump the current Major version
build       Will bump the current Build metadata
minor       Will bump the current Minor version
patch       Will bump the current Patch version
pre         Will bump the current Prerelease

<br/>

#### bump major
If our current version is `0.1.0`:
```
$ semver bump minor --> 1.0.0
```

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For example if our current version is `1.2.3-alpha`:

```
$ semver bump major --> 2.0.0
```

Usage:
```semver bump major```

<br/>

#### bump minor
If our current version is `0.1.0`:

```
$ semver bump minor --> 0.2.0
```

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For example if our current version is `1.2.3-alpha`:
```
$ semver bump minor --> 1.3.0
```
Usage:
```semver bump minor```

<br/>

#### bump patch
If our current version is `0.1.0`:

```
$ semver bump patch --> 0.1.1
```

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For example if our current version is `1.2.3-alpha`:

```$ semver bump patch --> 1.2.4```

Usage:
```semver bump patch```

<br/>

#### bump pre
`bump pre` means "create the next prerelease cut within the current prerelease stream."

If there is no prerelease yet, it starts one at `alpha.1`:

```text
$ semver bump pre
1.0.1 --> 1.0.1-alpha.1
```

If there is already a prerelease with a numeric tail, it increments that tail and keeps the current stage:

```text
$ semver bump pre
1.0.1-alpha.1 --> 1.0.1-alpha.2

$ semver bump pre
1.0.1-beta.1 --> 1.0.1-beta.2

$ semver bump pre
1.0.1-rc.2 --> 1.0.1-rc.3
```

If build metadata exists, `bump pre` removes it because a new prerelease cut is being created:

```text
$ semver bump pre
1.0.1-beta.2+ci.7 --> 1.0.1-beta.3
```

`bump pre` does not promote stages. To move between prerelease stages, use `set pre` explicitly:

```text
$ semver set pre beta.1
1.0.1-alpha.2 --> 1.0.1-beta.1

$ semver set pre rc.1
1.0.1-beta.4 --> 1.0.1-rc.1
```

Usage:
```semver bump pre```

<br/>

#### bump build
`bump build` increments existing build metadata. It does not invent a build label when one does not already exist.

```text
$ semver bump build
1.0.1+build.7 --> 1.0.1+build.8

$ semver bump build
1.0.1-rc.1+ci.7 --> 1.0.1-rc.1+ci.8
```

If no build metadata exists, `bump build` returns an error. Use `set build` first.

```text
$ semver set build ci.1
1.0.1 --> 1.0.1+ci.1
```

Usage:
```semver bump build```

---

### Set
By itself (with no subcommand) the set command will set the version to the passed in argument.  For example if our current version is 1.2.3:
$semver version 4.5.6 --> 4.5.6
$semver versiion 1.0.0-beta+exp.sha.5114f85 --> 1.0.0-beta+exp.sha.5114f85

Usage:
```
semver set [new-version]
semver set [command]
```

Available Commands:
build       Set Version Build information
pre         Set Version Pre Release information

<br/>

#### build
Will set the build on a version.  For example if the current version is 1.2.3:

```$ semver set build mybuild-123 --> 1.2.3+mybuild-123```

If no build string argument is given it will set the build to the short version of the current git HEAD hash.
This is equivalent to setting the build to the output of:

```$ git rev-parse --short HEAD```

For example if the current version is 1.2.3:

```$ semver set build --> 1.2.3+b113571 ```(if that was the current hash)

Build metadata is useful when you want extra release information that should not change SemVer precedence, such as:

- CI run numbers
- git SHAs
- packaging revisions
- internal build identifiers

Examples:

```text
$ semver set build ci.42
1.0.1 --> 1.0.1+ci.42

$ semver set build sha.abc1234
1.0.1-rc.1 --> 1.0.1-rc.1+sha.abc1234
```

Because build metadata does not affect SemVer precedence, these have the same precedence:

```text
1.0.1+ci.42
1.0.1+ci.99
```

Usage:
```semver set build [(optional) build value]```

<br/>

#### pre 
Will set the pre-release on a version.  For example if the current version is 1.2.3:

``$ semver set pre alpha-123 --> 1.2.3-alpha-123``

Prerelease identifiers are useful when you are cutting versions that are not yet final and should sort before the final release:

- `alpha` for earliest internal or limited testing
- `beta` for broader testing and stabilization
- `rc` for release candidates

Examples:

```text
$ semver set pre alpha.1
1.0.1 --> 1.0.1-alpha.1

$ semver set pre beta.1
1.0.1-alpha.3 --> 1.0.1-beta.1

$ semver set pre rc.1
1.0.1-beta.2 --> 1.0.1-rc.1
```

Prerelease precedence is lower than the final release, and the common progression is:

```text
alpha --> beta --> rc --> final
```

For example:

```text
1.0.1-alpha.1 < 1.0.1-alpha.2 < 1.0.1-beta.1 < 1.0.1-rc.1 < 1.0.1
```

Stage changes are explicit with `set pre`. They are not automatic in `bump pre`.

Common workflows:

```text
$ semver bump patch
1.0.0 --> 1.0.1

$ semver bump pre
1.0.1 --> 1.0.1-alpha.1

$ semver bump pre
1.0.1-alpha.1 --> 1.0.1-alpha.2

$ semver set pre beta.1
1.0.1-alpha.2 --> 1.0.1-beta.1

$ semver bump pre
1.0.1-beta.1 --> 1.0.1-beta.2

$ semver set pre rc.1
1.0.1-beta.2 --> 1.0.1-rc.1

$ semver set 1.0.1
1.0.1-rc.1 --> 1.0.1
```

Usage:
```semver set pre [(optional) pre-release value]```

Typical prerelease progression examples:

```text
v1.0.0-alpha.1
v1.0.0-alpha.2
v1.0.0-alpha.3
v1.0.0-beta.1
v1.0.0-rc.1
v1.0.0-rc.2
v1.0.0
```

The same progression is also valid without the leading `v`:

```text
1.0.0-alpha.1
1.0.0-alpha.2
1.0.0-alpha.3
1.0.0-beta.1
1.0.0-rc.1
1.0.0-rc.2
1.0.0
```

---
### Version
Prints the current version in the chosen format. For example if the current version is 1.2.3 Format options:

```
$ semver version --> 1.2.3 (string is default)
$ semver version -f string --> 1.2.3
$ semver version -f json -->
{
"Major": 2,
"Minor": 0,
"Patch": 0,
"PreRelease": "",
"Build": ""
}

$ semver version -f pretty --> {Major: 2, Minor: 0, Patch: 0, PreRelease: "", Build: ""}
```

Pretty differs form json in that pretty is a pretty print of the underlying Version struct and is technically not valid json.

Usage:
``semver version [flags]``

Flags:
-f, --format string   Print Format [string | json | pretty] (default "string")
