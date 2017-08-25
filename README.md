# stmts 

This is a utility package to create a concurrency safe collection of sqlx named statements that can be consumed in 3 stages. Registration, preparation, and execution. 

[![Go Report Card](https://goreportcard.com/badge/github.com/whitecypher/stmts)](https://goreportcard.com/report/github.com/whitecypher/stmts)

## how to

```go
package main

import (
  "fmt"
  "database/sql"
  
  "github.com/jmoiron/sqlx"
  "github.com/whitecypher/stmts"
)

// storing queries in constants will make it easier to retrieve the queries from the collection later. The query is the key.
const (
  qGetSomething = `SELECT * FROM mytable WHERE id = :id`
  qUpdateSomething = `UPDATE mytable SET my_name_is = :name WHERE id = :id`
  qDeleteSomething = `DELETE FROM mytable WHERE id = :id`
)

func main() { 
  // using sqlx for database operations
  var db *sqlx.DB
  
  // TODO: resolve the connection to the database
  
  // collection has no dependencies so a simple declaration is sufficient.
  var named stmts.Named
  
  // add queries
  named.Add(
    qGetSomething,
    qUpdateSomething,
    qDeleteSomething,
  )
  
  // prepare the statements
  err := named.Prepare(db)
  if err != nil {
    // TODO: handle the error however you like
  }
  // or use
  named.MustPrepare(db, func(format string, v ...interface{}) { 
    // TODO: handle the error however you like
  }) 
  
  // retrieve and execute a query
  var row struct {
    ID   uint64 `db:"id"`
    Name string `db:"name"`
  }
  err := named.Stmt(qGetSomething).Get(&row, struct {
    ID uint64 `db:"id"`
  }{
    ID: 1,
  })
  if err != nil && err != sql.ErrNoRows {
    // TODO: handle the error however you like
  }
}
```
