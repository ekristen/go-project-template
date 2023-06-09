# Go Project Template

This is an opinionated go project template to use as a starting point for new projects.

## Features

- Builds with [GoReleaser](https://goreleaser.com)
  - Automated with GitHub Actions
  - Signed with Cosign (providing you generate a private key)
- Builds with Docker
  - While designed to use goreleaser, you can still just run `docker build`
- Opinionated Layout
  - Never use `internal/` folder 
  - Everything is under `pkg/` folder
- Automatic Dependency Management with [Renovate](https://github.com/renovatebot/renovate)
- Automatic Releases with [Release Drafter](https://github.com/release-drafter/release-drafter)
- Documentation with Material for MkDocs
- API Server Example
  - Uses Gorilla Mux (yes it's been archived, still the best option)
- Stubbed out Go Tests
  - They are not comprehensive

## Building

The following will build binaries in snapshot order.

```console
goreleaser --clean --snapshot
```

## Configure

1. Rename Repository
2. Generate Cosign Keys
3. Update `.goreleaser.yml`
4. Update `main.go`  

### Signing

1. Create a password
   - Recommend exporting in environment as `COSIGN_PASSWORD`
2. Generate cosign keys
3. Create GitHub Action Secrets
   - `COSIGN_KEY` -> populate with cosign.key value
   - `COSIGN_PASSWORD` -> populate with password from step 1

#### Generate Keypair

```console
cosign generate-key-pair
```

## Documentation

The project is built to have the documentation right along side the code in the `docs/` directory.

This README.md is currently copied to `docs/index.md` and the documentation is automatically published to `gh-pages`
for the repository.