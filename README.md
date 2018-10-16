# Check Firewalled

Mini HTTP service that takes a JSON with server information and check
 if it's accessible from the application.

In our use-case, we want the tested server to **not** be accessible from this application.

An unaddressable server is considered a success.

## Env

| Key        | Requirement           | Explanation                                                                                   |
|------------|-----------------------|-----------------------------------------------------------------------------------------------|
| SHARED_KEY | mandatory             | Secret shared between main app and this one. (Needs to be sent in the header *Authorization*) |
| APP_PORT   | optional (default 80) | Port used for the application                                                                 |

## Data
```json
{
  "server": "example.com",
  "port": 25,
  "mx": false
}
```


| Key    | Explanation                                                                    |
|--------|--------------------------------------------------------------------------------|
| server | Server to check                                                                |
| port   | Port to use to attempt the connection                                          |
| mx     | Instead of resolving the IP, resolve the MX of the server first then check IPs |
## Result
### Success
```json
{
    "request": {
        "server": "example.com",
        "port": 25,
        "mx": false
    },
    "success": true,
    "details": [
        {
            "server": "example.com",
            "ip": "0.0.0.0",
            "result": true
        },
        {
            "server": "example.com",
            "ip": "0.0.0.1",
            "result": true
        },
    ]
}
```


### Failure

#### All servers
```json
{
    "request": {
        "server": "example.com",
        "port": 25,
        "mx": false
    },
    "success": false,
    "details": [
        {
            "server": "example.com",
            "ip": "0.0.0.0",
            "result": false
        },
        {
            "server": "example.com",
            "ip": "0.0.0.1",
            "result": false
        },
    ]
}
```

#### One servers
```json
{
    "request": {
        "server": "example.com",
        "port": 25,
        "mx": false
    },
    "success": false,
    "details": [
        {
            "server": "example.com",
            "ip": "0.0.0.0",
            "result": false
        },
        {
            "server": "example.com",
            "ip": "0.0.0.1",
            "result": true
        },
    ]
}
```