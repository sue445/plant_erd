## Unreleased
[full changelog](http://github.com/sue445/plant_erd/compare/v0.4.4...master)

## v0.4.4
[full changelog](http://github.com/sue445/plant_erd/compare/v0.4.3...v0.4.4)

* [Oracle] Fix for oracle open cursor and foreign key issue
  * https://github.com/sue445/plant_erd/pull/225
  * https://github.com/sue445/plant_erd/issues/224
* Upgrade to Go 1.21 :rocket:
  * https://github.com/sue445/plant_erd/pull/226
* Update dependencies

## v0.4.3
[full changelog](http://github.com/sue445/plant_erd/compare/v0.4.2...v0.4.3)

* [Sqlite3] Fixed. panic: interface conversion when foreign key `to` column is `nil`
  * https://github.com/sue445/plant_erd/pull/211
  * https://github.com/sue445/plant_erd/issues/210
* Upgrade to Go 1.20 :rocket:
  * https://github.com/sue445/plant_erd/pull/208
* Update dependencies

## v0.4.2
[full changelog](http://github.com/sue445/plant_erd/compare/v0.4.1...v0.4.2)

* Fix mermaid attribute error and display "not null"
  * https://github.com/sue445/plant_erd/pull/204
* Upgrade to Go 1.19 :rocket:
  * https://github.com/sue445/plant_erd/pull/194
* Update dependencies

## v0.4.1
[full changelog](http://github.com/sue445/plant_erd/compare/v0.4.0...v0.4.1)

* [PostgreSQL] Fixed foreign keys not included in public schema couldn't be retrieved
  * https://github.com/sue445/plant_erd/pull/182
  * https://github.com/sue445/plant_erd/issues/179
* Upgrade to Go 1.18
  * https://github.com/sue445/plant_erd/pull/176
* Support PostgreSQL 14
  * https://github.com/sue445/plant_erd/pull/177
* Update dependencies

## v0.4.0
[full changelog](http://github.com/sue445/plant_erd/compare/v0.3.0...v0.4.0)

* Upgrade to Go 1.17
  * https://github.com/sue445/plant_erd/pull/165
* Support mermaid :mermaid:
  * https://github.com/sue445/plant_erd/pull/170
* Add `--skip-table` argument to `plant_erd-oracle`
  * https://github.com/sue445/plant_erd/pull/171
* Update dependencies

## v0.3.0
[full changelog](http://github.com/sue445/plant_erd/compare/v0.2.4...v0.3.0)

* Add SkipTable patterns parameter
  * https://github.com/sue445/plant_erd/pull/161
* Update dependencies

## v0.2.4
[full changelog](http://github.com/sue445/plant_erd/compare/v0.2.3...v0.2.4)

* Fixed. `plant_erd sqlite3` doesn't work
  * https://github.com/sue445/plant_erd/pull/157
  * https://github.com/sue445/plant_erd/issues/155

## v0.2.3
[full changelog](http://github.com/sue445/plant_erd/compare/v0.2.2...v0.2.3)

* Exclude views on MySQL command
  * https://github.com/sue445/plant_erd/pull/158
  * https://github.com/sue445/plant_erd/issues/156

## v0.2.2
[full changelog](http://github.com/sue445/plant_erd/compare/v0.2.1...v0.2.2)

* Fixed. `version 'GLIBC_2.28' not found error` on Ubuntu 18.04
  * https://github.com/sue445/plant_erd/pull/153
  * https://github.com/sue445/plant_erd/issues/152
* Update dependencies

## v0.2.1
[full changelog](http://github.com/sue445/plant_erd/compare/v0.2.0...v0.2.1)

* Add darwin/arm64 (a.k.a. Apple M1) binary
  * https://github.com/sue445/plant_erd/pull/138
* Upgrade to Go 1.16
  * https://github.com/sue445/plant_erd/pull/137
* Update dependencies

## v0.2.0
[full changelog](http://github.com/sue445/plant_erd/compare/v0.1.1...v0.2.0)

* Support oracle
  * https://github.com/sue445/plant_erd/pull/50
* Update dependencies
  * https://github.com/sue445/plant_erd/pull/51
  * https://github.com/sue445/plant_erd/pull/56

## v0.1.1
[full changelog](http://github.com/sue445/plant_erd/compare/v0.1.0...v0.1.1)

* Resolved. doesn't work darwin bin with sqlite3
  * https://github.com/sue445/plant_erd/pull/49
  * https://github.com/sue445/plant_erd/issues/48

## v0.1.0
* Initial release
