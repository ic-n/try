# Try

The image of a baseball gopher catcher
Go package that enhances Go's error handling capabilities in microservices with complex partitioning into multiple layers.

Revisiting my old idea: [github.com/kiselev-nikolay/try](https://github.com/kiselev-nikolay/try)

This project suggested to be used as an inspiration only.

## Usage

Libruary implements way to handle errors and panics in classic try-catch constructions.

> `tc.Try` will run for first error.

```go
tc, cancel := try.New(context.Background())
defer cancel()

tc.Try(func() error {
    return fmt.Errorf("failed")
})
tc.Try(func() error {
    panic("this block will not be called")
})
tc.Catch(func(err error) {
    fmt.Print(err) // failed
})
```

> `tc.Try` also handles panics.

```go
tc, cancel := try.New(context.Background())
defer cancel()

tc.Try(func() error {
    panic("hello")
})
tc.Catch(func(err error) {
    fmt.Print(err) // panic: hello
})
```

> `tc.Try` will dispatch known errors.

```go
var anErr = errors.New("an error")

tc, cancel := try.New(context.Background())
defer cancel()

tc.Try(func() error {
    return anErr
})
tc.CatchError(anErr, func(err error) {
    fmt.Print(err) // an error
})
tc.Catch(func(other error) {
    // will not be executed
})
```

> `tc.Try` can ignore known errors.

```go
var anErr = errors.New("an error")

tc, cancel := try.New(context.Background())
defer cancel()

tc.Try(func() error {
    return anErr
})
tc.PassError(anErr)
tc.Catch(func(other error) {
    // will not be executed
})
```

__Complete example:__

```go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/ic-n/try"
	_ "modernc.org/sqlite"
)

func main() {
	var ( // assembly style variable declarations, let's go
		tc, cancel = try.New(context.Background())
		f          *os.File
		db         *sql.DB
		data       struct{ ID int }
	)
	defer cancel()

	tc.Try(func() (err error) {
		f, err = os.Open("object.json")
		return
	})
	tc.Try(func() (err error) {
		db, err = sql.Open("sqlite", "file::memory:")
		return
	})
	tc.Try(func() error {
		r, err := db.QueryContext(tc, "SELECT * FROM dataset;")
		if err != nil {
			return err
		}

		return r.Scan(data)
	})
	tc.Try(func() error {
		return json.NewEncoder(f).Encode(data)
	})

	tc.PassError(os.ErrExist)
	tc.Catch(func(err error) {
		log.Fatal(err)
	})
}
```
