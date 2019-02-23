# A Simple Forward Proxy

Rather than handling requests by itself, this proxy forwards requests to external proxy. For testing and default configuration, it uses [Tinyproxy](https://github.com/tinyproxy/tinyproxy), but it should works with other similar proxy server as well.

## Configuration

Configuration saved in `$HOME/.fwdproxy/config.conf` file using JSON format as follows,

```json
{
	"port": "8000",
	"proxy": {
		"address": "127.0.0.1",
		"port": "8888",
		"username": "test",
		"password": "testpassword"
	}
}
```

* **port**. Server port to listen HTTP request.
* **proxy**. External proxy configuration.
  * **address**. Address of external proxy.
  * **port**. Port of external proxy.
  * **username**. `Proxy-Authorization`'s user
  * **password**. `Proxy-Authorization`'s password.

## Default User
By default, access requires `Proxy-Authorization`. Credential informations kept in `$HOME/.fwdproxy/users.db` using [bolt](https://github.com/etcd-io/bbolt) database. On first use, a default user created as follows,

* Username: `admin`
* Password: `4dm1n`

More user can be added directly to database file.

## Log
Log messages is written to `$HOME/.fwdproxy/output.log`.

## License

This project is licensed under the BSD 2 License - see the [LICENSE](LICENSE) file for details.

