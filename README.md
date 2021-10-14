# PlantERD
ERD exporter with [PlantUML](https://plantuml.com/) format

[![Build Status](https://github.com/sue445/plant_erd/workflows/test/badge.svg?branch=master)](https://github.com/sue445/plant_erd/actions?query=workflow%3Atest)
[![Build Status](https://github.com/sue445/plant_erd/workflows/build/badge.svg?branch=master)](https://github.com/sue445/plant_erd/actions?query=workflow%3Abuild)
[![Coverage Status](https://coveralls.io/repos/github/sue445/plant_erd/badge.svg?branch=master)](https://coveralls.io/github/sue445/plant_erd?branch=master)
[![Maintainability](https://api.codeclimate.com/v1/badges/0a9432880ae3f992cc65/maintainability)](https://codeclimate.com/github/sue445/plant_erd/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/sue445/plant_erd)](https://goreportcard.com/report/github.com/sue445/plant_erd)

## Example
```bash
$ ./plant_erd sqlite3 --database /path/to/test_db.sqlite3

entity articles {
  * id : integer
  --
  * user_id : integer
  --
  index_user_id_on_articles (user_id)
}

entity users {
  * id : integer
  --
  name : text
}

articles }-- users
```

![example](./img/example.svg)

## Features
* Output ERD from real database
* Output ERD to stdout or file
* Output only tables within a certain distance adjacent to each other with foreign keys from a specific table

## Supports
* SQLite3
* MySQL: 5.6, 5.7 and 8
* PostgreSQL: 9, 10, 11, 12 and 13
* Oracle

## Setup
Download latest binary from https://github.com/sue445/plant_erd/releases and `chmod 755`

* `plant_erd` : for SQLite3, MySQL and PostgreSQL
* `plant_erd-oracle` : for Oracle

### Setup for `plant_erd-oracle`
`plant_erd-oracle` requires Basic Package or Basic Light Package in [Oracle Instant Client](https://www.oracle.com/database/technologies/instant-client.html)

#### Example (Linux)
```bash
mkdir -p /opt/oracle
wget --quiet --tries=0 https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip
unzip -q instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip -d /opt/oracle
export LD_LIBRARY_PATH=/opt/oracle/instantclient_19_3

# for Ubuntu
apt-get update
apt-get install -y libaio1
```

#### Example (Mac)
See https://github.com/kubo/ruby-oci8/blob/master/docs/install-on-osx.md and install `instantclient-basic` or `instantclient-basiclite`

#### Example (Windows)
1. Go to https://www.oracle.com/database/technologies/instant-client/winx64-64-downloads.html
2. Download `instantclient-basic-windows.x64-19.5.0.0.0dbru.zip` or ` instantclient-basiclite-windows.x64-19.5.0.0.0dbru.zip`
3. Extract zip
4. Move `plant_erd-oracle` to same directory as `oci.dll`

## Usage
### SQLite3
```
$ ./plant_erd sqlite3 --help
NAME:
   plant_erd sqlite3 - Generate ERD from sqlite3

USAGE:
   plant_erd sqlite3 [command options] [arguments...]

OPTIONS:
   -d DISTANCE, --distance DISTANCE  Output only tables within a certain DISTANCE adjacent to each other with foreign keys from a specific table (default: 0)
   --database DATABASE               SQLite3 DATABASE file
   -f FILE, --file FILE              FILE for output (default: stdout)
   -i, --skip-index                  Whether don't print index to ERD
   -s value, --skip-table value      Skip generating table by using regex patterns
   -t TABLE, --table TABLE           Output only tables within a certain distance adjacent to each other with foreign keys from a specific TABLE
```

### MySQL
```bash
$ ./plant_erd mysql --help
NAME:
   plant_erd mysql - Generate ERD from mysql

USAGE:
   plant_erd mysql [command options] [arguments...]

OPTIONS:
   --collation COLLATION             MySQL COLLATION (default: "utf8_general_ci")
   -d DISTANCE, --distance DISTANCE  Output only tables within a certain DISTANCE adjacent to each other with foreign keys from a specific table (default: 0)
   --database DATABASE               MySQL DATABASE name
   -f FILE, --file FILE              FILE for output (default: stdout)
   --host HOST                       MySQL HOST (default: "localhost")
   -i, --skip-index                  Whether don't print index to ERD
   --password PASSWORD               MySQL PASSWORD [$MYSQL_PASSWORD]
   --port PORT                       MySQL PORT (default: 3306)
   -s value, --skip-table value      Skip generating table by using regex patterns
   -t TABLE, --table TABLE           Output only tables within a certain distance adjacent to each other with foreign keys from a specific TABLE
   --user USER                       MySQL USER (default: "root")
```

### PostgreSQL
```bash
$ ./plant_erd postgresql --help
NAME:
   plant_erd postgresql - Generate ERD from PostgreSQL

USAGE:
   plant_erd postgresql [command options] [arguments...]

OPTIONS:
   -d DISTANCE, --distance DISTANCE  Output only tables within a certain DISTANCE adjacent to each other with foreign keys from a specific table (default: 0)
   --database DATABASE               PostgreSQL DATABASE name
   -f FILE, --file FILE              FILE for output (default: stdout)
   --host HOST                       PostgreSQL HOST (default: "localhost")
   -i, --skip-index                  Whether don't print index to ERD
   --password PASSWORD               PostgreSQL PASSWORD [$POSTGRES_PASSWORD]
   --port PORT                       PostgreSQL PORT (default: 5432)
   -s value, --skip-table value      Skip generating table by using regex patterns
   --sslmode SSLMODE                 PostgreSQL SSLMODE. c.f. https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-PARAMKEYWORDS (default: "disable")
   -t TABLE, --table TABLE           Output only tables within a certain distance adjacent to each other with foreign keys from a specific TABLE
   --user USER                       PostgreSQL USER
```

### Oracle
```bash
$ ./plant_erd-oracle --help
NAME:
   plant_erd-oracle - ERD exporter with PlantUML format (for oracle)

USAGE:
   plant_erd-oracle [global options] command [command options] [arguments...]

VERSION:
   vX.X.X (build. xxxxxxx)

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -f FILE, --file FILE              FILE for output (default: stdout)
   -t TABLE, --table TABLE           Output only tables within a certain distance adjacent to each other with foreign keys from a specific TABLE
   -d DISTANCE, --distance DISTANCE  Output only tables within a certain DISTANCE adjacent to each other with foreign keys from a specific table (default: 0)
   -i, --skip-index                  Whether don't print index to ERD
   --user USER                       Oracle USER
   --password PASSWORD               Oracle PASSWORD [$ORACLE_PASSWORD]
   --host HOST                       Oracle HOST (default: "localhost")
   --port PORT                       Oracle PORT (default: 1521)
   --service SERVICE                 Oracle SERVICE name
   --help, -h                        show help
   --version, -v                     print the version
```

## About `--table` and `--distance`
When `--table` and `--distance` are passed, output only tables within a certain distance adjacent to each other with foreign keys from a specific table.

### Example 1: Output all tables
```bash
$ ./plant_erd sqlite3
```

![example all](img/example-all.svg)

### Example 2: Output only tables within a distance of 1 from the articles
```bash
$ ./plant_erd sqlite3 --table articles --distance 1
```

![example distance 1 from articles](img/example-distance-1-from-articles.svg)

## Testing
### with all databases
Run test in container

```bash
docker-compose up --build --abort-on-container-exit
```

### with only SQLite3
Run test on local

```bash
make test
```

## License
The program is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).

But `plant_erd-oracle` contains [Oracle Instant Client](https://www.oracle.com/database/technologies/instant-client.html).
Oracle Instant Client is under [OTN License](https://www.oracle.com/downloads/licenses/distribution-license.html).
