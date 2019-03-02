## Rigby

[Rigby](https://rigby.space) is a micro-blogging platform made in Go. 

To set up your own instance of rigby, you need a few things:

#### Config file
An example `config.json` file - HTTPS is optional:

```json
{
	{
	"database": {
		"username": "luke",
		"password": "mypassword1",
		"database": "rigby"
	}, 
	"server": {
		"port": 3000,
		"origin": "localhost" 
	},
	"https": {
		"certificate": "/path/to/cert",
		"key": "/path/to/key"
	}
}
```

#### MySQL server
You'll need a MySQL server that has `utf8mb4` encoding enabled. Currently there are no schema files for the MySQL server setup. Those should be coming soon. 