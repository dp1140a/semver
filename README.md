# semver

Yet another semantic versioning tool written in golang that follows the [Semver v2 Spec](https://semver.org/)
I wrote this cause it follows how I like to use these tools: less flags and more subcommands that make sense.

## Getting Started

### Install
If you already have golang installed you can install by running the command:

```sh
go install github.com/dp1140a/semver/cmd/semver@latest
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
minor       Will bump the current Minor version
patch       Will bump the current Patch version

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

---

### Set
By itself (with no subcommand) the set command will set the version to the passed in argument.  For example if our current version is 1.2.3:
$semver version 4.5.6 --> 4.5.6
$semver versiion 1.0.0-beta+exp.sha.5114f85 --> 1.0.0-beta+exp.sha.5114f85

Usage:
```s
emver set [new-version]
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

Usage:
```semver set build [(optional) build value]```

<br/>

#### pre 
Will set the pre-release on a version.  For example if the current version is 1.2.3:

``$ semver set pre alpha-123 --> 1.2.3-alpha-123``

If no pre string argument is given it will set the pre-release accordingly:
If no pre-release value will set to alpha.  For example if the current version is 1.2.3

```$ semver set pre --> 1.2.3-alpha```

If pre-release value is alpha will set to beta.  For example if the current version is 1.2.3-alpha

```$ semver set pre --> 1.2.3-beta```

If pre-release value is beta will set to rc1.0.  For example if the current version is 1.2.3-beta

```$ semver set pre --> 1.2.3-rc1.0```

NOTE: Setting the pre-release value WILL delete the current build value since pre-release is a higher precedence.

Usage:
```semver set pre [(optional) pre-release value]```

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
