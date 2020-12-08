# githubsecret

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/jefflinse/githubsecret/Main%20CI)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/jefflinse/githubsecret)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jefflinse/githubsecret)
![GitHub](https://img.shields.io/github/license/jefflinse/githubsecret)

A Go package for encrypting GitHub secrets. It uses hashing and encryption APIs from [golang.org/x/crypto](https://golang.org/x/crypto) and does not require libsodium C bindings.

## Usage

```shell
$ go get github.com/jefflinse/githubsecret
```

```go
package main

import (
    "github.com/jefflinse/githubsecret"
)

func main() {
    // 1. obtain your repository's public key:
    // https://docs.github.com/en/free-pro-team@latest/rest/reference/actions#get-a-repository-public-key

    // 2. encrypt your secret
    encrypted, err := githubsecret.Encrypt(repoPublicKey, content)

    // 3. store the encrypted secret
    // https://docs.github.com/en/free-pro-team@latest/rest/reference/actions#create-or-update-a-repository-secret
}
```

## Example

The examples directory contains a CLI app called `putsecret` for storing GitHub secrets. It demonstrates the entire workflow of obtaining a repository's public key, using it to encrypt a secret, and storing the secret for use in GitHub Actions.

Clone, build, and run it with no arguments to view its usage:

```shell
$ git clone https://github.com/jefflinse/githubsecret

$ cd githubsecret/examples/putsecret

$ go build

$ ./putsecret
```

Define `GITHUB_USERNAME` and `GITHUB_TOKEN` in your environment with your GitHub username and personal access token, respectively. Your access token must have sufficient privileges to read the repository's public key and to update secrets.

Pass the owner, repository, secret name (key), and secret value as command line arguments.

```shell
$ ./putsecret owner repo secret_id "secret value"
```

Go to the Secrets page in your repository's settings and you should see your secret listed.

## Contributing

1. Make your changes in a new branch.
1. Test your changes locally and make sure all UTs pass. Add or update tests to ensure all code is covered.
1. Push your branch and create a pull request.
1. Label your PR with one of `major release`, `minor release`, or `patch release` depending on the nature of your changes.
1. After approval and merging, [pr-semver-bump](https://github.com/jefflinse/pr-semver-bump) will automatically create a new release.
