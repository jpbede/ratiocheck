# ratiocheck
Microservice to check image to content ratio of HTML pages

## Requirements
When running outside the docker container a Chrome installation is required. 

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
### Shell

You can run a check by issuing following command:

```shell
ratiocheck check <url to check>
```

### REST API

The REST API accepts and returns JSON. Successful requests return `200 OK`
with a ratio result:

```json
{
  "content_area": 1591200,
  "image_area": 948332,
  "ratio": 59.598541980894915
}
```

The response fields are:

| Field | Description |
| --- | --- |
| `content_area` | Total rendered document area in square pixels. |
| `image_area` | Combined rendered area of all `<img>` elements in square pixels. |
| `ratio` | Percentage of the rendered document area covered by images. For example, `59.6` means images cover about 59.6% of the page. |

Invalid JSON or missing required fields return `400 Bad Request`:

```json
{
  "error": "invalid json"
}
```

Pages that cannot be loaded or measured return `500 Internal Server Error` with
an `error` message.

#### HTML
Do an HTTP POST call to `/html` endpoint with following JSON body:

```json
{
  "html": "<your html>"
}
```

The service writes the HTML to a temporary local file and measures the rendered
page in headless Chrome.

#### URL
Do an HTTP POST call to `/url` endpoint with following JSON body:

```json
{
  "url": "<your url>"
}
```
