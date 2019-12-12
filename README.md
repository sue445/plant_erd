# PlantERD
ERD exporter with [PlantUML](https://plantuml.com/) format

[![Build Status](https://github.com/sue445/plant_erd/workflows/test/badge.svg?branch=master)](https://github.com/sue445/plant_erd/actions?query=workflow%3Atest)
[![Coverage Status](https://coveralls.io/repos/github/sue445/plant_erd/badge.svg?branch=master)](https://coveralls.io/github/sue445/plant_erd?branch=master)

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
* PostgreSQL: 9, 10, 11 and 12

## Usage
Download latest binary from https://github.com/sue445/plant_erd/releases and `chmod 755`

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
   --sslmode SSLMODE                 PostgreSQL SSLMODE. c.f. https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-PARAMKEYWORDS (default: "disable")
   -t TABLE, --table TABLE           Output only tables within a certain distance adjacent to each other with foreign keys from a specific TABLE
   --user USER                       PostgreSQL USER
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

### without MySQL and PostgreSQL
Run test on local

```bash
make test
```
