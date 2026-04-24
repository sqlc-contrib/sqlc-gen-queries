# Changelog

## [0.2.0](https://github.com/sqlc-contrib/sqlc-gen-queries/compare/v0.1.0...v0.2.0) (2026-04-24)


### Features

* rewrite list queries with comment-marker placeholders ([4026f85](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/4026f85eeb08cca7775a337c07dda8e8fca033a7))

## [0.1.0](https://github.com/sqlc-contrib/sqlc-gen-queries/compare/v0.0.5...v0.1.0) (2026-04-16)


### Features

* distinguish foreign key indexes from other non-unique indexes ([00775d5](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/00775d5b34b4f35b1bc808275c5355c1bba4c9a8))

## [0.0.5](https://github.com/sqlc-contrib/sqlc-gen-queries/compare/v0.0.4...v0.0.5) (2026-04-16)


### Bug Fixes

* update nix vendorHash for changed go dependencies ([0efbf12](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/0efbf12a74820d3c20a908f4e92dc592fa894a31))

## [0.0.4](https://github.com/sqlc-contrib/sqlc-gen-queries/compare/v0.0.3...v0.0.4) (2026-04-16)


### Features

* make list queries included by default ([d7104a8](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/d7104a83bb53103e608692c2d0a5e624e1faa0f6))


### Bug Fixes

* use full git history in build job for octocov to push coverage files ([4a64ffe](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/4a64ffef8cecc6a6605aed0e08234ae805da7711))

## [0.0.3](https://github.com/sqlc-contrib/sqlc-gen-queries/compare/v0.0.2...v0.0.3) (2026-03-17)


### Bug Fixes

* squeeze consecutive blank lines in generated output ([be78ef4](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/be78ef41d726deae69ba92711ccf90259296cdc5))

## [0.0.2](https://github.com/sqlc-contrib/sqlc-gen-queries/compare/v0.0.1...v0.0.2) (2026-03-13)


### Features

* add SQL query template generation with CRUD operations ([de0b1d8](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/de0b1d869938fa191459a5355803fed5705bddbb))
* **config:** support multiple config file locations ([4328412](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/4328412d4cf2aed6c86adf30f407dc40687cee20))
* implement CLI application with config and catalog loading ([12e573f](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/12e573ffbbcac0a884f121987d7aefba411c965f))
* **inflect:** add singular and plural forms for "quota" ([3259507](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/325950769af180ec3797cfa6be630f6a479812fa))
* rework config to use queries allowlist instead of skip_queries blocklist ([0b29e6a](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/0b29e6a9955bdb3813e54fa8006ec8996ccc07e9))
* **sqlc:** add GetNonPrimaryKeyColumns method to filter primary keys from updates ([fdb05fd](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/fdb05fd43da0ea587252f89a5d0f7234137ed4c5))
* **sqlc:** add skip_queries configuration for selective query generation ([8879fd7](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/8879fd758cce67a669e8e56b45d0210569f30908))


### Bug Fixes

* **ci:** add report.path so octocov tracks generated files for push ([13968fe](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/13968fe79f1c58f316209cf60dad3fb5d0ca97e4))
* **ci:** reorganize octocov output into .github/octocov/ directory ([929ce65](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/929ce65c0692cb9765bb1eb1f6098311acc7e07f))
* **ci:** restore octocov badge config and lower coverage threshold to 80% ([c041571](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/c041571a9bd524604399fa8641b37a441707d35c))
* **template:** add blank line separators between generated query blocks ([feb407b](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/feb407b1dba9ea4e626e0b6d631f926d1cc79c76))
* **test:** add missing config_test_exclude.yaml fixture ([#17](https://github.com/sqlc-contrib/sqlc-gen-queries/issues/17)) ([2579c1e](https://github.com/sqlc-contrib/sqlc-gen-queries/commit/2579c1e89c30f6b9d345287949fbda22795eb4bf))
