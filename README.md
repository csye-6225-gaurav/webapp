# webapp

## Requirements

- Go (version 1.23.1)
- PostgreSQL(15)

## Installation

- clone the repoository:

```bash
git clone https://github.com/csye-6225-gaurav/webapp.git
cd webapp
```
- Install Go dependencies:

```bash
go mod tidy
```
Before running the application, create a `.env` file in the root directory of the project. It should contain the following details:

```bash
DB_Host=
DB_Port=
DB_Pass=
DB_User=
DB_Name=
DB_SSLMode=
APP_Port=
```
- Run the application:
```bash
go run main.go
```