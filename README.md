# ratiocheck
Microservice to check image to content ration of HTML pages

**Please be nice to me, the project is work in progress :)**

## Installing

### Binaries
You find pre-compiled binaries and packages for the most common OS under the [releases](https://github.com/jpbede/ratiocheck/releases).

### Docker
```shell
docker run -p 3000:3000 ghcr.io/jpbede/ratiocheck:latest
```

or via `docker-compose.yml`

```yaml
version: "3"
services:
  ratiocheck:
    image: ghcr.io/jpbede/ratiocheck:latest
    ports:
      - 3000:3000
```

### macOS
Simply use `homebrew` (https://brew.sh/)

To install `ratiocheck` use following command `brew install jpbede/tap/ratiocheck`

## Using it

Do an HTTP POST call to `/ratio` endpoint with following JSON body:

```json
{
  "url": "<your url>"
}
```

As a result you will get the ratio image area to content area, the size of the content area and the size of the image area:

```json
{
  "content_area": 1591200,
  "image_area": 948332,
  "ratio": 59.598541980894915
}
```