# Changelog

## 0.3.1 (2023-11-27)


### Features

* implement traefik github oauth plugin ([d3be0a5](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/d3be0a5831ad83a7e8ceab47e0d6216902755313))
* implement traefik github oauth server app ([7a7acdf](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/7a7acdf7f9822dee89225b3a17b3ac732bef5c94))
* **middleware:** add log ([789e4cf](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/789e4cf0209aa13cd1aff5302a679686e63fcf29))
* **middleware:** message when api secret key is invalid ([6138346](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/61383468b262150387da2f7a9598d8984a01dbde))
* **middleware:** use `github.com/dghubble/sling` as http client ([81f461f](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/81f461fb35ed3fc5aa9d3441aec6c3a29e8f3db4))
* **server:** add log ([48cf8ea](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/48cf8ea367d4c033918c2a4c2ca15148da1b32a8))
* **server:** return request error message in json ([4c1eac9](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/4c1eac941db36e701f97d32335406b57bfafa860))
* **server:** use chi-router instead of gin server ([e862713](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/e8627136aa97344d8d28d5cad9c2c012066f6ce2))
* set no cache headers ([316878f](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/316878f0d3f2e8fa04a8eb6697c3a924eecd66c5))
* **traefik-plugin:** ensure the correct is always sent and add local dockerfile for testing ([0fa2208](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/0fa22086a48e15ad865d03c4134b4a73ed216d7c))
* update vendor packages and add chi router ([d021f58](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/d021f58498674f928c295de6be98cc535952b3a8))
* use httpin to get the correct form/json fields ([0f2f511](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/0f2f511960e07277f8427b67bf960523c6999d63))
* **vendor:** update vendor packages ([0b5975d](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/0b5975dd0864f77aa17e892e6d7418142f2f2552))


### Bug Fixes

* change package name to avoid conflicts ([f7bef93](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/f7bef9329d5bb5615b239bdc4ae5c270c29ee0e0))
* ignore changes to the dist folder ([2f91dd8](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/2f91dd88499bf3f5b1a796f6b7ffc86122751587))
* **middleware:** redirect only on get requests ([61af42c](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/61af42ceb3917f44a0ef0aee5c2678fac670e164))
* **server:** fix incorrect use of context ([788a2b0](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/788a2b0514bed2ae13252f60e104e9d3a4aa1ff2))
* **server:** Fix logger middleware log fields ([3ccd7e3](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/3ccd7e38015495f2a91c31e2342d299baf86ae25))


### Continuous Integration

* build multi-platform Docker images using goreleaser ([fda884c](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/fda884c3d6887dad11c05620d287d8e3aa9efe41)), closes [#22](https://github.com/luizfonseca/traefik-github-oauth-plugin/issues/22)

## [0.5.0](https://github.com/luizfonseca/traefik-github-oauth-plugin/compare/v0.4.2...v0.5.0) (2023-11-26)


### Features

* use httpin to get the correct form/json fields ([0f2f511](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/0f2f511960e07277f8427b67bf960523c6999d63))

## [0.4.2](https://github.com/luizfonseca/traefik-github-oauth-plugin/compare/v0.4.1...v0.4.2) (2023-11-26)


### Bug Fixes

* change package name to avoid conflicts ([f7bef93](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/f7bef9329d5bb5615b239bdc4ae5c270c29ee0e0))

## [0.4.1](https://github.com/luizfonseca/traefik-github-oauth-plugin/compare/v0.4.0...v0.4.1) (2023-11-26)


### Bug Fixes

* ignore changes to the dist folder ([2f91dd8](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/2f91dd88499bf3f5b1a796f6b7ffc86122751587))

## [0.4.0](https://github.com/luizfonseca/traefik-github-oauth-plugin/compare/v0.3.1...v0.4.0) (2023-11-26)


### Features

* **server:** use chi-router instead of gin server ([e862713](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/e8627136aa97344d8d28d5cad9c2c012066f6ce2))
* update vendor packages and add chi router ([d021f58](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/d021f58498674f928c295de6be98cc535952b3a8))
* **vendor:** update vendor packages ([0b5975d](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/0b5975dd0864f77aa17e892e6d7418142f2f2552))

## [0.3.1](https://github.com/luizfonseca/traefik-github-oauth-plugin/compare/v0.3.1...v0.3.1) (2023-11-26)


### Features

* implement traefik github oauth plugin ([d3be0a5](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/d3be0a5831ad83a7e8ceab47e0d6216902755313))
* implement traefik github oauth server app ([7a7acdf](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/7a7acdf7f9822dee89225b3a17b3ac732bef5c94))
* **middleware:** add log ([789e4cf](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/789e4cf0209aa13cd1aff5302a679686e63fcf29))
* **middleware:** message when api secret key is invalid ([6138346](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/61383468b262150387da2f7a9598d8984a01dbde))
* **middleware:** use `github.com/dghubble/sling` as http client ([81f461f](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/81f461fb35ed3fc5aa9d3441aec6c3a29e8f3db4))
* **server:** add log ([48cf8ea](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/48cf8ea367d4c033918c2a4c2ca15148da1b32a8))
* **server:** return request error message in json ([4c1eac9](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/4c1eac941db36e701f97d32335406b57bfafa860))
* set no cache headers ([316878f](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/316878f0d3f2e8fa04a8eb6697c3a924eecd66c5))


### Bug Fixes

* **middleware:** redirect only on get requests ([61af42c](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/61af42ceb3917f44a0ef0aee5c2678fac670e164))
* **server:** fix incorrect use of context ([788a2b0](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/788a2b0514bed2ae13252f60e104e9d3a4aa1ff2))
* **server:** Fix logger middleware log fields ([3ccd7e3](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/3ccd7e38015495f2a91c31e2342d299baf86ae25))


### Continuous Integration

* build multi-platform Docker images using goreleaser ([fda884c](https://github.com/luizfonseca/traefik-github-oauth-plugin/commit/fda884c3d6887dad11c05620d287d8e3aa9efe41)), closes [#22](https://github.com/luizfonseca/traefik-github-oauth-plugin/issues/22)

## [0.3.1](https://github.com/MuXiu1997/traefik-github-oauth-plugin/compare/v0.3.0...v0.3.1) (2023-11-15)


### Continuous Integration

* build multi-platform Docker images using goreleaser ([fda884c](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/fda884c3d6887dad11c05620d287d8e3aa9efe41)), closes [#22](https://github.com/MuXiu1997/traefik-github-oauth-plugin/issues/22)

## [0.3.0](https://github.com/MuXiu1997/traefik-github-oauth-plugin/compare/v0.2.2...v0.3.0) (2023-02-04)


### Features

* set no cache headers ([316878f](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/316878f0d3f2e8fa04a8eb6697c3a924eecd66c5))

## [0.2.2](https://github.com/MuXiu1997/traefik-github-oauth-plugin/compare/v0.2.1...v0.2.2) (2023-01-30)


### Bug Fixes

* **server:** fix incorrect use of context ([788a2b0](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/788a2b0514bed2ae13252f60e104e9d3a4aa1ff2))

## [0.2.1](https://github.com/MuXiu1997/traefik-github-oauth-plugin/compare/v0.2.0...v0.2.1) (2023-01-27)


### Bug Fixes

* **server:** Fix logger middleware log fields ([3ccd7e3](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/3ccd7e38015495f2a91c31e2342d299baf86ae25))

## [0.2.0](https://github.com/MuXiu1997/traefik-github-oauth-plugin/compare/v0.1.1...v0.2.0) (2023-01-27)


### Features

* **middleware:** add log ([789e4cf](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/789e4cf0209aa13cd1aff5302a679686e63fcf29))
* **server:** add log ([48cf8ea](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/48cf8ea367d4c033918c2a4c2ca15148da1b32a8))


### Bug Fixes

* **middleware:** redirect only on get requests ([61af42c](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/61af42ceb3917f44a0ef0aee5c2678fac670e164))

## 0.1.1 (2023-01-26)


### Features

* implement traefik github oauth plugin ([d3be0a5](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/d3be0a5831ad83a7e8ceab47e0d6216902755313))
* implement traefik github oauth server app ([7a7acdf](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/7a7acdf7f9822dee89225b3a17b3ac732bef5c94))
* **middleware:** message when api secret key is invalid ([6138346](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/61383468b262150387da2f7a9598d8984a01dbde))
* **middleware:** use `github.com/dghubble/sling` as http client ([81f461f](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/81f461fb35ed3fc5aa9d3441aec6c3a29e8f3db4))
* **server:** return request error message in json ([4c1eac9](https://github.com/MuXiu1997/traefik-github-oauth-plugin/commit/4c1eac941db36e701f97d32335406b57bfafa860))
