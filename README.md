## Rigby

[Rigby](rigby.space) is a micro-blogging platform made in Go. 

To set up your own instance of rigby, you need a few things:

#### Config file
An example `config.json` file:

```json
	{
	"database": {
		"username": "luke",
		"password": "mypassword1",
		"database": "rigby"
	}, 
	"server": {
		"port": 3000,
		"origin": "localhost" 
	}
}
```

#### MySQL server
You'll need a MySQL server that has `utf8mb4` encoding enabled. Currently there are no schema files for the MySQL server setup. Those should be coming soon. 

#### Golang compiler
You'll also need the [Go Programming Language](https://golang.org/doc/install) tools installed in order to compile the project.