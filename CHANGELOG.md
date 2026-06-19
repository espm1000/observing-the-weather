# Changelog

## [0.1.2](https://github.com/espm1000/observing-the-weather/compare/v0.1.1...v0.1.2) (2026-06-19)


### Bug Fixes

* **build:** update image push condition to trigger on release events ([#58](https://github.com/espm1000/observing-the-weather/issues/58)) ([f6fb345](https://github.com/espm1000/observing-the-weather/commit/f6fb3456b34db0313dd141e4bc5d58571b435020))

## [0.1.1](https://github.com/espm1000/observing-the-weather/compare/v0.1.0...v0.1.1) (2026-06-19)


### Bug Fixes

* **build:** update image tagging to use release tag instead of timestamp ([#56](https://github.com/espm1000/observing-the-weather/issues/56)) ([c6ebb31](https://github.com/espm1000/observing-the-weather/commit/c6ebb31c46ba008a61cef7b9b1d5c1b69d780253))

## 0.1.0 (2026-06-19)


### Features

* add data structures for weather observations and forecasts ([528c6e0](https://github.com/espm1000/observing-the-weather/commit/528c6e08a038bbdfca6d883f2d96f8efd467e579))
* implement weather data retrieval and forecasting functionality ([77230ff](https://github.com/espm1000/observing-the-weather/commit/77230ff19ef4f3898d6bb68a2c989cee96895b11))
* initialize project structure and add initial documentation ([b681c15](https://github.com/espm1000/observing-the-weather/commit/b681c156bd6187efc8b2f5ce26c1f946887934d9))
* **logging:** enhance logging and reporting for weather data ([#41](https://github.com/espm1000/observing-the-weather/issues/41)) ([1183877](https://github.com/espm1000/observing-the-weather/commit/1183877dadbcddff301566a295ccb811e6066357))
* **output:** update CSV output handling and add Docker support ([#39](https://github.com/espm1000/observing-the-weather/issues/39)) ([ac2045a](https://github.com/espm1000/observing-the-weather/commit/ac2045a2cadfcc445b50996e81b52052482fd3f3))
* **release-please:** add automated release management workflow ([#51](https://github.com/espm1000/observing-the-weather/issues/51)) ([dbfab63](https://github.com/espm1000/observing-the-weather/commit/dbfab638cb5a2b611e2bd601b2c62948bc0b7dbc))
* **report:** implement CSV generation for current weather data ([#38](https://github.com/espm1000/observing-the-weather/issues/38)) ([393550e](https://github.com/espm1000/observing-the-weather/commit/393550e3d86db673bf850fb314db0597db7038df))
* **weather:** implement weather data retrieval and forecasting functionality ([c8ac1a9](https://github.com/espm1000/observing-the-weather/commit/c8ac1a98a650cbe44715c45b07b4e87888cfe3fd))
* **workflow:** add GitHub Actions workflow for building and pushing Docker image ([#40](https://github.com/espm1000/observing-the-weather/issues/40)) ([e5b4103](https://github.com/espm1000/observing-the-weather/commit/e5b41036903fb6a585428dce0d01cb804cc4633f))
* **workflow:** add grype integration and security jobs ([#44](https://github.com/espm1000/observing-the-weather/issues/44)) ([124d3bd](https://github.com/espm1000/observing-the-weather/commit/124d3bd07c537232362ddf95debb4a7c0a10cc3b))


### Bug Fixes

* improve error handling for non-200 response in GetCurrentWeather ([937443c](https://github.com/espm1000/observing-the-weather/commit/937443ced280431f77d6a74053c6093dba2606a3))
* **logger:** ensure log directory is created before logging setup and update image tag format in GitHub Actions workflow ([#42](https://github.com/espm1000/observing-the-weather/issues/42)) ([143fc06](https://github.com/espm1000/observing-the-weather/commit/143fc061324af3f0e3bfcdbd687c68e41f0d972e))
* **release-please:** set separate-pull-requests to false and restore package config ([#53](https://github.com/espm1000/observing-the-weather/issues/53)) ([00ef580](https://github.com/espm1000/observing-the-weather/commit/00ef5804a85860659f33f2e34cc48650eead14d4))
* **workflow:** correct condition for skipping chore branches in build job ([#48](https://github.com/espm1000/observing-the-weather/issues/48)) ([7a274e5](https://github.com/espm1000/observing-the-weather/commit/7a274e56193072c7d1727ffce1099b2262b61479))


### Miscellaneous Chores

* manual release ([dbb56e1](https://github.com/espm1000/observing-the-weather/commit/dbb56e19f329b00a5b07b7ec34187e89c69887f3))
