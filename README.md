# Mattrax
Open Source MDM (Mobile Device Management) System.

**Currently only Windows MDM** is being supported but in the future I would like to support other MDM protocols such as Apple (IOS & MacOS), Android and ChromeOS.

## Project Status

This project is under heavy development. You can use it but expect to do some serious debugging and **don't expose to the internet** as all security mechanisms have not been implemented.

## Running

This project requires an external [PostgreSQL](https://www.postgresql.org/) database. The `sql/schema.sql` and `sql/deploy.sql` should be run on a blank database to configure it. You should change the `deploy.sql` to fit your deployment settings. Then start the Go binary (with arguments `--db "postgres://localhost/Mattrax" --domain mdm.example.com`) and your MDM server will be working.

## Developing

This project uses [sqlc](https://github.com/kyleconroy/sqlc) so the command `sqlc generate` is used to generate the `internal/db` package from `sql/queries.sql`.