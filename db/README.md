We use golang-migrate for DB schema management. 'migrations/' directory contains files for it.

When running in docker-compose *.up.sql files are added to MySQL container and provisioned on startup.

When running in Kubernetes the 'dataservice' image runs migrations as post-install hook.

For running migrations manually, install 'migrate' utility from https://github.com/golang-migrate/migrate/releases and run it like that from the current directory:

```shell
migrate -database "mysql://blog:blogadmin123@tcp(10.152.183.90:3306)/blog" -path=migrations/ up
```
