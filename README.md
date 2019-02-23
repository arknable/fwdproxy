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


## License

This project is licensed under the BSD 2 License - see the [LICENSE](LICENSE) file for details.
