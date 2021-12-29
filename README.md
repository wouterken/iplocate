# IPLocate Go
This is a basic self-hostable server written in go-lang that can accept incoming HTTP requests
and return basic location information about the requesting IP address.

Credits:
* https://github.com/oschwald/maxminddb-golang
* https://dev.maxmind.com/geoip/geolite2-free-geolocation-data

## Usage
You can run the go-process as standalone, or build the included docker image.
Both have a dependency on a file named "GeoLite2-City.mmdb" in the root directory.

```
docker build . -t iplocate
docker run -p --rm -it 3000:3000 --name iplocate iplocate
```