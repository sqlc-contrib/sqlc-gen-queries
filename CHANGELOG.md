# Changelog

## 1.0.0 (2026-03-13)


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
