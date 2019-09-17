# MAGENTO CONSUMER SERVICE


## Requirements
1. Instal go
2. [Install soda](https://gobuffalo.io/en/docs/db/toolbox/) (to create database and running migration)


## Running Service
- `go install`
- `go get -u`
- `go run main.go` to run the server

## Debugging
- Please to update `REQUEST_DEBUG` to true, to debug request traffic
- If you want to log all information to File, please to update `IS_LOG_FILE` to true.

## Database
- `soda g config` to create `database.yml` configuration
- `soda create -e development` (to create database development) [more](https://gobuffalo.io/en/docs/db/toolbox/)
- `soda generate fizz name_of_migration` to create new migrations.
- `soda drop -e development` (to drop or delete database) [more](https://gobuffalo.io/en/docs/db/toolbox/)
- migration up : `soda migrate -p database up` (`database` is folder which is where migrations folder laid) [more](https://gobuffalo.io/en/docs/db/migrations/)
- migration down : `soda migrate -p database down -s {number of database want to down}`
