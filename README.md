# Mock JSON Server

## How to use it

```sh
./server-<operating_system> --file <path to json fixture> --port <port number>
```

### OPTIONS

--port :default 3000 \
--file : it is required

### Example JSON File

```json
{
    "routes": [
      {
        "path": "/users",
        "method": "GET",
        "data": "json://example/data/users.json" //relative path to the executable
      },
      {
        "path": "/settings",
        "method": "POST",
        "data": {"theme": "dark"}
      },
      {
        "path": "/users",
        "method": "GET",
        "data": [
            {"id": 1, "name": "John Doe"},
            {"id": 2, "name": "Jane Smith"}
            ]   
        }
    ]
}
```

### Download the latest release from the [Releases](https://github.com/agniswarm/json-mock-server/releases)
