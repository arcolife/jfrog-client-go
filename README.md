# Bintray CLI
Bintray CLI

Bintray CLI provides a command line interface for invoking actions on Bintray.

## Getting Started

### Downloading the executables from Bintray

[TO DO]

### Building the command line executable

If you prefer, you may instead build the client in go.

#### Setup GO on your machine

* Make sure you have a working Go environment. [See the install instructions](http://golang.org/doc/install).
* Navigate to the directory where you want to create the *bintray-cli-go* project.
* Set the value of the GOPATH environment variable to the full path of this directory.

#### Download Bintray CLI from GitHub

Run the following command to create the *bintray-cli-go* project:
```console
$ go get github.com/JFrogDev/bintray-cli-go/...
```

Navigate to the following directory
```console
$ cd $GOPATH/bin
```
### Usage

You can copy the *bt* executable to any location on your file-system as long as you add it to your *PATH* environment variable,
so that you can access it from any path.

#### Command syntax

```console
$ bt command-name global-options command-options arguments
```

The sections below specify the available commands, their respective options and additional arguments that may be needed.
*bt* should be followed by a command name (for example, download-ver), a list of options (for example, --repo=my-bintray-repo)
and the list of arguments for the command.

#### Global options

Global options are used for all commands.
```console
   --user            [Optional] Bintray username. It can be also set using the BINTRAY_USER environment variable. If not set, the subject sent as part of the command argument is used for authentication.
   --key             [Mandatory] Bintray API key. It can be also set using the BINTRAY_KEY environment variable.
   --api-url         [Default: https://api.bintray.com] Bintray API URL. It can be also set using the BINTRAY_API_URL environment variable.
   --download-url    [Default: https://dl.bintray.com] Bintray download server URL. It can be also set using the BINTRAY_DOWNLOAD_URL environment variable.
```

#### Commands list
- [upload (u)](#upload)
- [download-file (dlf)](#download-file)
- [download-ver (dlv)](#download-ver)
- [package-show (ps)](#package-show)
- [package-create (pc)](#package-create)
- [package-update (pc)](#package-update)
- [package-delete (pd)](#package-delete)
- [version-show (vs)](#version-show)
- [version-create (vc)](#version-create)
- [version-update (vu)](#version-update)
- [version-delete (vd)](#version-delete)
- [version-publish (vd)](#version-publish)
- [entitlement-keys (ent-keys)](#entitlement-keys)
- [entitlements (ent)](#entitlements)
- [gpg-sign-file (gsf)](#gpg-sign-file)
- [gpg-sign-ver (gsv)](#gpg-sign-ver)

<a name="upload"/>
#### The *upload* (u) command

##### Function
Used to upload files to Bintray

##### Command options
```console
   --flat                        [Default: true]   If not set to true, and the upload path ends with a slash, artifacts are uploaded according to their file system hierarchy.
   --recursive                   [Default: true]   Set to false if you do not wish to collect artifacts in sub-directories to be uploaded to Bintray.
   --regexp                      [Default: false]  Set to true to use a regular expression instead of wildcards expression to collect artifacts to upload.
   --dry-run                     [Default: false]  Set to true to disable communication with Bintray.
```
If the Bintray Package to which you're uploading does not exist, the CLI will try to create it.
Please send the following command options for the package creation.
```console
   --pkg-desc                    [Optional]        Package description.
   --pkg-labels                  [Optional]        Package lables in the form of "lable11","lable2"...
   --pkg-licenses                [Mandatory]       Package licenses in the form of "Apache-2.0","GPL-3.0"...
   --pkg-cust-licenses           [Optional]        Package custom licenses in the form of "my-license-1","my-license-2"...
   --pkg-vcs-url                 [Mandatory]       Package VCS URL.
   --pkg-website-url             [Optional]        Package web site URL.
   --pkg-i-tracker-url           [Optional]        Package Issues Tracker URL.
   --pkg-github-repo             [Optional]        Package Github repository.
   --pkg-github-rel-notes        [Optional]        Github release notes file.
   --pkg-pub-dn                  [Default: false]  Public download numbers.
   --pkg-pub-stats               [Default: false]  Public statistics
```
If the Package Version to which you're uploading does not exist, the CLI will try to create it.
Please send the following command options for the version creation.   
```console
   --ver-desc                    [Optional]   Version description.
   --ver-vcs-tag                 [Optional]   VCS tag.
   --ver-released                [Optional]   Release date in ISO8601 format (yyyy-MM-dd'T'HH:mm:ss.SSSZ).
   --ver-github-rel-notes        [Optional]   Github release notes file.
   --ver-github-tag-rel-notes    [Optional]   Set to true if you wish to use a Github tag release notes.
```

##### Arguments
* The first argument is the local file-system path to the artifacts to be uploaded to Bintray.
The path can include a single file or multiple artifacts, by using the * wildcard.
**Important:** If the path is provided as a regular expression (with the --regexp=true option) then
the first regular expression appearing as part of the argument must be enclosed in parenthesis.

* The second argument is the upload path in Bintray.
The path can include symbols in the form of {1}, {2}, ...
These symbols are replaced with the sections enclosed with parenthesis in the first argument.

##### Examples

Upload all files located under *dir/sub-dir*, with names that start with *frog*, to the root path under version *1.0* 
of the *froggy-package* package 
```console
bt u "dir/sub-dir/frog*" "my-org/swamp-repo/froggy-package/1.0/" --user=my-user --key=my-api-key
```

Upload all files located under *dir/sub-dir*, with names that start with *frog* to the /frog-files folder, under version *1.0* 
of the *froggy-package* package 
```console
bt u "dir/sub-dir/frog*" "my-org/swamp-repo/froggy-package/1.0/frog-files/" --user=my-user --key=my-api-key
```

Upload all files located under *dir/sub-dir* with names that start with *frog* to the root path under version *1.0*, 
while adding the *-up* suffix to their names in Bintray.  
```console
bt u "dir/sub-dir/(frog*)" "my-org/swamp-repo/froggy-package/1.0/{1}-up" --user=my-user --key=my-api-key
```
<a name="download-file"/>
#### The *download-file* (dlf) command

##### Function
Used to download a specific file from Bintray.

##### Command options
This command has no command options. It uses however the global options.

##### Arguments
The command expects one argument in the form of *subject/repository/package/version/path*.

##### Examples
```console
bt download-file my-org/swamp-repo/froggy-package/1.0/com/jfrog/bintray/crazy-frog.zip --user=my-user --key=my-api-key
bt dlf my-org/swamp-repo/froggy-package/1.0/com/jfrog/bintray/crazy-frog.zip --user=my-user --key=my-api-key
```

<a name="download-ver"/>
#### The *download-ver* (dlv) command

##### Function
Used to download the files of a specific version from Bintray.

##### Command options
This command has no command options. It uses however the global options.

##### Arguments
The command expects one argument in the form of *subject/repository/package/version*.

##### Examples
```console
bt download-ver my-org/swamp-repo/froggy-package/1.0 --user=my-user --key=my-api-key
bt dlv my-org/swamp-repo/froggy-package/1.0 --user=my-user --key=my-api-key
```

<a name="package-show"/>
#### The *package-show* (ps) command

##### Function
Used for showing package details.

##### Command options
This command has no command options. It uses however the global options.

##### Arguments
The command expects one argument in the form of *subject/repository/package.

##### Examples
Show package *super-frog-package* 
```console
bt ps my-org/swamp-repo/super-frog-package 
```

<a name="package-create"/>
#### The *package-create* (pc) command

##### Function
Used for creating a package in Bintray

##### Command options
The command uses the global options, in addition to the following command options.
```console
   --desc               [Optional]        Package description.
   --labels             [Optional]        Package lables in the form of "lable11","lable2"...
   --licenses           [Mandatory]       Package licenses in the form of "Apache-2.0","GPL-3.0"...
   --cust-licenses      [Optional]        Package custom licenses in the form of "my-license-1","my-license-2"...
   --vcs-url            [Mandatory]       Package VCS URL.
   --website-url        [Optional]        Package web site URL.
   --i-tracker-url      [Optional]        Package Issues Tracker URL.
   --github-repo        [Optional]        Package Github repository.
   --github-rel-notes   [Optional]        Github release notes file.
   --pub-dn             [Default: false]  Public download numbers.
   --pub-stats          [Default: false]  Public statistics
```

##### Arguments
The command expects one argument in the form of *subject/repository/package.

##### Examples
Create the *super-frog-package* package 
```console
bt pc my-org/swamp-repo/super-frog-package --licenses=Apache-2.0,GPL-3.0 --vcs-url=http://github.com/jfrogdev/coolfrog.git 
```

<a name="package-update"/>
#### The *package-update* (pu) command

##### Function
Used for updating package details in Bintray

##### Command options
The command uses the same option as the *package-create* command

##### Arguments
The command expects one argument in the form of *subject/repository/package.

##### Examples
Create the *super-frog-package* package 
```console
bt pu my-org/swamp-repo/super-frog-package --labels=label1,label2,label3 
```

<a name="package-delete"/>
#### The *package-delete* (pd) command

##### Function
Used for deleting a packages in Bintray

##### Command options
The command uses the global options, in addition to the following command option.
```console
   --q      [Default: false]       Set to true to skip the delete confirmation message.
```

##### Arguments
The command expects one argument in the form of *subject/repository/package.

##### Examples
Delete the *froger-package* package 
```console
bt pc my-org/swamp-repo/froger-package --licenses=Apache-2.0,GPL-3.0 --vcs-url=http://github.com/jfrogdev/coolfrog.git 
```



<a name="version-show"/>
#### The *version-show* (vs) command

##### Function
Used for showing version details.

##### Command options
This command has no command options. It uses however the global options.

##### Arguments
The command expects one argument in one of the forms
* *subject/repository/package* to show the latest published version.
* *subject/repository/package/version* to show the specified version.

##### Examples
Show version 1.0.0 of package *super-frog-package* 
```console
bt vs my-org/swamp-repo/super-frog-package/1.0 
```

Show the latest published version of package *super-frog-package* 
```console
bt vs my-org/swamp-repo/super-frog-package 
```

<a name="version-create"/>
#### The *version-create* (vc) command

##### Function
Used for creating a version in Bintray

##### Command options
The command uses the global options, in addition to the following command options.
```console
   --desc                    [Optional]   Version description.
   --vcs-tag                 [Optional]   VCS tag.
   --released                [Optional]   Release date in ISO8601 format (yyyy-MM-dd'T'HH:mm:ss.SSSZ).
   --github-rel-notes        [Optional]   Github release notes file.
   --github-tag-rel-notes    [Optional]   Set to true if you wish to use a Github tag release notes.
```

##### Arguments
The command expects one argument in the form of *subject/repository/package/version.

##### Examples
Create version 1.0.0 in package *super-frog-package* 
```console
bt vc my-org/swamp-repo/super-frog-package/1.0.0 
```

<a name="version-update"/>
#### The *version-update* (vu) command

##### Function
Used for updating version details in Bintray

##### Command options
The command uses the same option as the *version-create* command

##### Arguments
The command expects one argument in the form of *subject/repository/package/version.

##### Examples
Update the labels of version 1.0.0 in package *super-frog-package* 
```console
bt vu my-org/swamp-repo/super-frog-package/1.0.0 --labels=jump,jumping,frog
```

<a name="version-delete"/>
#### The *version-delete* (vd) command

##### Function
Used for deleting a version in Bintray

##### Command options
The command uses the global options, in addition to the following command option.
```console
   --q      [Default: false]       Set to true to skip the delete confirmation message.
```

##### Arguments
The command expects one argument in the form of *subject/repository/package/version.

##### Examples
Create version 1.0.0 in package *super-frog-package* 
```console
bt vd my-org/swamp-repo/super-frog-package/1.0.0 
```

<a name="version-publish"/>
#### The *version-publish* (vp) command

##### Function
Used for publishing a version in Bintray

##### Command options
This command has no command options. It uses however the global options.

##### Arguments
The command expects one argument in the form of *subject/repository/package/version.

##### Examples
Publish version 1.0.0 in package *super-frog-package* 
```console
bt vp my-org/swamp-repo/super-frog-package/1.0.0 
```

<a name="entitlement-keys"/>
#### The *entitlement-keys* (ent-keys) command

##### Function
Used for managing Entitlement Download Keys.

##### Command options
The command uses the global options, in addition to the following command options.
```console
   --org             Bintray organization.
   --expiry          Download Key expiry (required for 'bt ent-keys show/create/update/delete'
   --ex-check-url    Used for Download Key creation and update. You can optionally provide an existence check directive, in the form of a callback URL, to verify whether the source identity of the Download Key still exists.
   --ex-check-cache  Used for Download Key creation and update. You can optionally provide the period in seconds for the callback URK results cache.
   --white-cidrs     Used for Download Key creation and update. Specifying white CIDRs in the form of 127.0.0.1/22,193.5.0.1/22 will allow access only for those IPs that exist in that address range.
   --black-cidrs     Used for Download Key creation and update. Specifying black CIDRs in the foem of 127.0.0.1/22,193.5.0.1/22 will block access for all IPs that exist in the specified range.
```

##### Arguments
* With no arguments, a list of all download keys is displayed.
* When sending the show, create, update or delete argument, it should be followed by an argument indicating the download key ID to show, create, update or delete the download key.

##### Examples
Show all Download Keys
```console
bt ent-keys
```
Create a Download Key
```console
bt ent-keys create key1 
bt ent-keys create key1 --expiry=7956915742000 
```

Show a specific Download Key
```console
bt ent-keys show key1
```

Update a Download Key
```console
bt ent-keys update key1 --ex-check-url=http://new-callback.com --white-cidrs=127.0.0.1/22,193.5.0.1/92 --black-cidrs=127.0.0.1/22,193.5.0.1/92
bt ent-keys update key1 --expiry=7956915752000
```

Delete a Download Key
```console
bt ent key delete key1
```

<a name="entitlements"/>
#### The *entitlements* (ent) command

##### Function
Used for managing Entitlements.

##### Command options
The command uses the global options, in addition to the following command options.
```console
   --id              Entitlement ID. Used for Entitlements update.
   --access          Entitlement access. Used for Entitlements creation and update.
   --keys            Used for Entitlements creation and update. List of Download Keys in the form of \"key1\",\"key2\"...
   --path            Entitlement path. Used for Entitlements creating and update.
```

##### Arguments
* When sending an argument in the form of subject/repo or subject/repo/package or subject/repo/package/version, a list of all entitles for the send repo, package or version is displayed.
* If the first argument sent is show, create, update or delete, it should be followed by an argument in the form of subject/repo or subject/repo/package or subject/repo/package/version to show, create, update or delete the entitlement.

##### Examples

Show all Entitlements of the swamp-repo repository.
```console
bt ent my-org/swamp-repo
```

Show all Entitlements of the green-frog package.
```console
bt ent my-org/swamp-repo/green-frog
```

Show all Entitlements of version 1.0 of the green-frog package.
```console
bt ent my-org/swamp-repo/green-frog/1.0
```

Create an Entitlement for the green-frog package, with rw access, the key1 and key2 Download Keys and the a/b/c path.
```console
bt ent create my-org/swamp-repo/green-frog --access=rw --keys=key1,key2 --path=a/b/c
```

Show a specific Entitlement on the swamp-repo repository.
```console
bt ent show my-org/swamp-repo --id=451433e7b3ec3f18110ba770c77b9a3cb5534cfc
```

Update the download keys and access of an Entitlement on the swamp-repo repository.
```console
bt ent update my-org/swamp-repo --id=451433e7b3ec3f18110ba770c77b9a3cb5534cfc --keys=key1,key2 --access=r
```

Delete an Entitlement on the my-org/swamp-repo.
```console 
bt ent delete my-org/swamp-repo --id=451433e7b3ec3f18110ba770c77b9a3cb5534cfc
```

<a name="sign-url"/>
#### The *sign-url* (su) command

##### Function
Used for Generating an anonymous, signed download URL with an expiry date.

##### Command options
The command uses the global options, in addition to the following command options.
```console
   --expiry            [Optional]    An expiry date for the URL, in Unix epoch time in milliseconds, after which the URL will be invalid. By default, expiry date will be 24 hours.
   --valid-for         [Optional]    The number of seconds since generation before the URL expires. Mutually exclusive with the --expiry option.
   --callback-id       [Optional]    An applicative identifier for the request. This identifier appears in download logs and is used in email and download webhook notifications.
   --callback-email    [Optional]    An email address to send mail to when a user has used the download URL. This requiers a callback_id. The callback-id will be included in the mail message.
   --callback-url      [Optional]    A webhook URL to call when a user has used the download URL.
   --callback-method   [Optional]    HTTP method to use for making the callback. Will use POST by default. Supported methods are: GET, POST, PUT and HEAD.
```

##### Arguments
The command expects one argument in the form of *subject/repository/file-path*.

##### Examples
Create a download URL for *froggy-file*, located under *froggy-folder* in the *swamp-repo* repository.
 ```console
bt us my-org/swamp-repo/froggy-folder/froggy-file
```

<a name="gpg-sign-file"/>
#### The *gpg-sign-file* (gsf) command

##### Function
GPG sign a file in Bintray.

##### Command options
The command uses the global options, in addition to the following command option.
```console
   --passphrase      [Optional]    GPG passphrase.
```

##### Arguments
The command expects one argument in the form of *subject/repository/file-path*.

##### Examples
GPG sign the *froggy-file* file, located under *froggy-folder* in the *swamp-repo* repository.
 ```console
bt gsf my-org/swamp-repo/froggy-folder/froggy-file
```
or with a passphrase
```console
bt gsf my-org/swamp-repo/froggy-folder/froggy-file --passphrase=gpgX***yH8eKw
```

<a name="gpg-sign-ver"/>
#### The *gpg-sign-ver* (gsv) command

##### Function
GPS sign all files of a specific version in Bintray.

##### Command options
The command uses the global options, in addition to the following command option.
```console
   --passphrase      [Optional]    GPG passphrase.
```

##### Arguments
The command expects one argument in the form of *subject/repository/package/version*.

##### Examples
GPG sign all files of version *1.0* of the *froggy-package* package in the *swamp-repo* repository.
 ```console
bt gsv my-org/swamp-repo/froggy-package/1.0
```
or with a passphrase
```console
bt gsv my-org/swamp-repo/froggy-package/1.0 --passphrase=gpgX***yH8eKw
```