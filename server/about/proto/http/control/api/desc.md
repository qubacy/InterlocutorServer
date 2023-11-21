# Http - Control 🤖

## Some details

### Error response

- *400*
```json
{
    "text": "<string>",
    "details": "<string>"
}
```

<!-- --------------------------------------------------------------------- -->

## POST /control/api/sign-in

### Request body
```json
{
    "login": "<string>",
    "password": "<string>"
}   
```

### Responses
- *200*
```json
{
    "access-token": "<jwt-string>"
}
```

<!-- --------------------------------------------------------------------- -->

## GET /control/api/admin

### Request headers
| Key | Value Type |
|-----|------------|
| Authorization | Bearer `jwt-string` | 

### Responses
- *200*
```json
{
    "admins": [
        {
            "idr": "<int>",
            "login": "<string>"
        },
        ...
    ]
}
```

## POST /control/api/admin

### Request headers
| Key | Value Type |
|-----|------------|
| Authorization | Bearer `jwt-string` | 

### Request body
```json
{
    "admin": {
        "login": "<string>",
        "password": "<string>",
    }
}
```

### Responses
- *200*
```json
{
    "idr": "<int>"
}

```

<!-- --------------------------------------------------------------------- -->

## GET /control/api/topic

### Request headers
| Key | Value Type |
|-----|------------|
| Authorization | Bearer `jwt-string` | 

### Responses
- *200*
```json
{
    "topics": [
        {
            "idr": "<int>",
            "lang": "<int>",
            "name": "<string>"
        },
        ...
    ]
}
```

## POST /control/api/topic
### Request headers
| Key | Value Type |
|-----|------------|
| Authorization | Bearer `jwt-string` | 

### Request body
```json
{
    "topics": [
        {
            "lang": "<int>",
            "name": "<string>"
        },
        ...
    ]
}
```

### Responses
- *200*