[![CircleCI](https://circleci.com/gh/HnH/qry/tree/master.svg?style=svg&circle-token=cd6ef5c602e0f89a80488349a1e4fbe034b8d717)](https://circleci.com/gh/HnH/qry/tree/master)
[![codecov](https://codecov.io/gh/HnH/qry/branch/master/graph/badge.svg)](https://codecov.io/gh/HnH/qry)
[![Go Report Card](https://goreportcard.com/badge/github.com/HnH/qry)](https://goreportcard.com/report/github.com/HnH/qry)

# About

**qry** is a general purpose library for storing your raw database queries in .sql files with all benefits of modern IDEs, instead of strings and constants in the code, and using them in an easy way inside your application with all the profit of compile time constants.

**qry** recursively loads all .sql files from a specified folder, parses them according to predefined rules and returns a reusable object, which is actually just a `map[string]string` with some sugar. Multiple queries inside a single file are separated with standard SQL comment syntax: `-- qry: QueryName`. A `QueryName` must match `[A-Za-z_]+`.

[gen](https://github.com/HnH/qry/tree/master/cmd/qry-gen) tool is used for automatic generation of constants for all user specified `query_names`.

# Installation

`go get -u github.com/HnH/qry/cmd/qry-gen/...` 

# Usage

Prepare sql files: `queries/one.sql`:

```sql
-- qry: InsertUser
INSERT INTO `users` (`name`) VALUES (?);

-- qry: GetUserById
SELECT * FROM `users` WHERE `user_id` = ?;
```

And the second one `queries/two.sql`:

```sql
-- qry: DeleteUsersByIds
DELETE FROM `users` WHERE `user_id` IN ({ids});
```

[Gen](https://github.com/HnH/qry/tree/master/cmd/qry-gen)erate constants: `qry-gen -dir=./queries -pkg=/path/to/your/go/pkg` Will produce `/path/to/your/go/pkg/qry.go` with:

```go
package pkg

const (
	// one.sql
	InsertUser  = "INSERT INTO `users` (`name`) VALUES (?);"
	GetUserById = "SELECT * FROM `users` WHERE `user_id` = ?;"

	// two.sql
	DeleteUsersByIds = "DELETE FROM `users` WHERE `user_id` IN ({ids});"
)
```

As a best practice include this qry-gen call in your source code with go:generate prefix: `//go:generate qry-gen -dir=./queries -pkg=/path/to/your/go/pkg` and just execute `go generate` before each build.
Now it's time to use **qry** inside your project:

```go
func main() {
	/**
	 * The most obvious way is to use generated constants in the source code
	 */
	 
	// INSERT INTO `users` (`name`) VALUES (?);
	println(pkg.InsertUser)
	
	// DELETE FROM `users` WHERE `user_id` IN (?,?,?);
	println(qry.Query(pkg.DeleteUsersByIds).Replace("{ids}", qry.In(3)))
	
	/**
	 * As an alternative you can manually parse .sql files in the directory and work with output
	 */
	if q, err := qry.Dir("/path/to/your/go/pkg/queries"); err != nil {
		log.Fatal(err)
	}

	// SELECT * FROM `users` WHERE `user_id` = ?;
	println(q["one.sql"]["GetUserById"])
  
	// DELETE FROM `users` WHERE `user_id` IN (?,?,?);
	println(q["two.sql"]["DeleteUsersByIds"].Replace("{ids}", qry.In(3)))
}
```
