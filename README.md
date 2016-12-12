# zero  [![Build Status](https://travis-ci.org/umayr/zero.svg?branch=master)](https://travis-ci.org/umayr/zero) [![Go Report Card](https://goreportcard.com/badge/github.com/umayr/zero)](https://goreportcard.com/report/github.com/umayr/zero)
> super tiny in-memory experimental store

#### Build

```
λ go get -u github.com/umayr/zero
λ cd $GOPATH/src/github.com/umayr/zero
λ make
```

This will create two binary files `zero-client` and `zero-server`. You can use as many client as you want, data will be shared among all clients as long as they're connect to same running server.

#### Supported Types

- `string`
- `number`
- `array<string|number>`

#### Commands

- `ADD <key> <value>`
- `SHOW <*|key>`
- `KEYS`
- `COUNT`
- `DEL <key>`
- `PUSH <key> <value>`
- `POP <key>`
- `EXIT`

#### Example
```
λ ./build/zero-client
> ADD foo hello
OK
> ADD bar world
OK
> SHOW *
+-----+-------+--------+---------------------+
| KEY | VALUE |  TYPE  |        TIME         |
+-----+-------+--------+---------------------+
| foo | hello | string | 12 Dec 16 05:43 PKT |
| bar | world | string | 12 Dec 16 05:43 PKT |
+-----+-------+--------+---------------------+

> KEYS
[foo bar]
> ADD n 1000
OK
> ADD x 19928300012000
OK
> SHOW *
+-----+----------------+--------+---------------------+
| KEY |     VALUE      |  TYPE  |        TIME         |
+-----+----------------+--------+---------------------+
| n   |           1000 | number | 12 Dec 16 05:43 PKT |
| x   | 19928300012000 | number | 12 Dec 16 05:44 PKT |
| foo | hello          | string | 12 Dec 16 05:43 PKT |
| bar | world          | string | 12 Dec 16 05:43 PKT |
+-----+----------------+--------+---------------------+

> ADD arr [1, 2, 3]
OK
> SHOW arr
[1 2 3]
> PUSH arr hello world
4
> SHOW arr
[1 2 3 hello world]
> SHOW *
+-----+-----------------------+--------+---------------------+
| KEY |         VALUE         |  TYPE  |        TIME         |
+-----+-----------------------+--------+---------------------+
| x   |        19928300012000 | number | 12 Dec 16 05:44 PKT |
| n   |                  1000 | number | 12 Dec 16 05:43 PKT |
| bar | world                 | string | 12 Dec 16 05:43 PKT |
| foo | hello                 | string | 12 Dec 16 05:43 PKT |
| arr | [1 2 3 hello world]   | array  | 12 Dec 16 05:44 PKT |
+-----+-----------------------+--------+---------------------+

> POP arr
hello world
> POP arr
 3
> POP arr
 2
> POP arr
1
> POP arr
err: array is already empty
> KEYS
[foo bar n x arr]
> DEL arr
OK
> COUNT
4
> DEL foo
OK
> DEL n
OK
> SHOW *
+-----+----------------+--------+---------------------+
| KEY |     VALUE      |  TYPE  |        TIME         |
+-----+----------------+--------+---------------------+
| x   | 19928300012000 | number | 12 Dec 16 05:44 PKT |
| bar | world          | string | 12 Dec 16 05:43 PKT |
+-----+----------------+--------+---------------------+

> 
```
#### License
MIT
