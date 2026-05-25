## Rigby

Rigby is a micro-blogging platform.

To set up your own instance of rigby, you need a few things:

#### Config file
An example `config.json` file:

```json
{
  "database": "rigby.db",
  "server": {
    "port": 3000,
    "origin": "localhost"
  },
  "rss": {
    "baseUrl": "http://localhost:3000"
  }
}
```

#### SQLite
Rigby uses SQLite for its database

#### Golang compiler
You'll also need the [Go Programming Language](https://golang.org/doc/install) tools installed in order to compile the project.
