# Changelog

## [1.2.0](https://github.com/olunusib/go-ci/compare/v1.1.3...v1.2.0) (2024-05-29)


### Features

* support logs ([c115c8a](https://github.com/olunusib/go-ci/commit/c115c8a58fd7a6040c5c780d77804c6a309a8de9))

## [1.1.3](https://github.com/olunusib/go-ci/compare/v1.1.2...v1.1.3) (2024-05-29)


### Bug Fixes

* only build image on release ([8ede2fa](https://github.com/olunusib/go-ci/commit/8ede2fa3c8440910201bf26f323f59a72fb16418))

## [1.1.2](https://github.com/olunusib/go-ci/compare/v1.1.1...v1.1.2) (2024-05-29)


### Bug Fixes

* get pipeline path after repo clone ([df49ad5](https://github.com/olunusib/go-ci/commit/df49ad5e881276f2e856eda5b8d70e1b40e594a9))

## [1.1.1](https://github.com/olunusib/go-ci/compare/v1.1.0...v1.1.1) (2024-05-28)


### Miscellaneous Chores

* release 1.1.1 ([61a557f](https://github.com/olunusib/go-ci/commit/61a557f6ae5eb5dcfe3a32f92fc67b8e7544ca47))

## [1.1.0](https://github.com/olunusib/go-ci/compare/v1.0.0...v1.1.0) (2024-05-28)


### Features

* add dockerignore file ([4c8650c](https://github.com/olunusib/go-ci/commit/4c8650c4c432a18480fdb199d048fcbb7b9ba7ef))
* add functionality for loading pipeline and running simple steps ([aa3fbf4](https://github.com/olunusib/go-ci/commit/aa3fbf47fb2831532652f9695715807cc8a38750))
* add release-please config ([bf37be0](https://github.com/olunusib/go-ci/commit/bf37be05c66d4d2393a5d38fdd9ab6c0d19d42e4))
* add support for env vars ([97cb89a](https://github.com/olunusib/go-ci/commit/97cb89acb434073c893c4d6e776669ef945308bb))
* add support for more platforms with QEMU in docker image ([41feab6](https://github.com/olunusib/go-ci/commit/41feab60108814575b18056567df6a2c0b2ae1c2))
* add web server and expose a webhook endpoint ([e05dc92](https://github.com/olunusib/go-ci/commit/e05dc9269ebfd0bf3e944e3ae1722d5e5139915a))
* enable commit statuses ([035320e](https://github.com/olunusib/go-ci/commit/035320ef90251dbd307631d61c42dc3dfc523a4c))
* use multistep build for docker image ([8867a1c](https://github.com/olunusib/go-ci/commit/8867a1c16dd94a0b14441cecd38f893717c10813))


### Bug Fixes

* add id-token permission to deploy workflow ([63c1ee7](https://github.com/olunusib/go-ci/commit/63c1ee76e2756d62cfb99350159ffa7349756f09))
* go 1.22 fails to build for armv7 ([be218e1](https://github.com/olunusib/go-ci/commit/be218e179377095e02b7a8e478b11bf35d721c41))
* temporarily disable artifact attestation ([e10678a](https://github.com/olunusib/go-ci/commit/e10678a1c7a474ad603901f2834045ca802c0455))
* temporarily disable artifact attestation ([2ccb155](https://github.com/olunusib/go-ci/commit/2ccb15520b76dcc0943e376601048108a0ddf844))
* use alpine base image instead ([e4ec6ab](https://github.com/olunusib/go-ci/commit/e4ec6ab7bc0e53a4f4dbec1665813ccef74a82fb))
* use PAT for release-please ([f4e9288](https://github.com/olunusib/go-ci/commit/f4e9288acff526afa3d512b8aadae3f1f81cb096))
