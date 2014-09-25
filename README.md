sqlc
----

[![Build Status](https://travis-ci.org/relops/sqlc.png?branch=master)](https://travis-ci.org/relops/sqlc)

`sqlc` generates SQL for you:
	
	foo := Table("foo")
	bar := Varchar(foo, "bar")
	baz := Varchar(foo, "baz")

	var db *db.DB // For integration with database/sql

	row, err := Select(bar).From(foo).Where(baz.Eq("quux")).QueryRow(db)

If you don't want to use `database/sql`, you don't have to. `String()` is an API call to just produce the SQL string that you use in any way that you want to:

	// Renders `SELECT foo.bar FROM foo WHERE foo.baz = ?`
	sql := Select(bar).From(foo).Where(baz.Eq("quux")).String()

Features
--------

* SELECT ... FROM ... WHERE
* GROUP BY
* ORDER BY
* INNER JOINS
* Data types:
  * VARCHAR
* Statement rendering
* Querying via database/sql


Status
------

Experimental - this is work in progress. Basically I'm trying to port [JOOQ][] to Go, but I don't know yet whether it will work.

[jooq]: http://jooq.org