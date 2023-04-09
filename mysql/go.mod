module github.com/mozgio/database/mysql

go 1.20

replace github.com/mozgio/database => ../

require (
	github.com/go-sql-driver/mysql v1.7.0
	github.com/mozgio/database v0.0.0-00010101000000-000000000000
)

require github.com/adlio/schema v1.3.3 // indirect
