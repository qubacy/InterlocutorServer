# Examples

## OPS: `Searching start`, `Searching game found`, `Chatting new message`, etc.

### Request
#### `Searching start`
```json
{
    "operation": 0,
    "body": {
        "profile": {
            "username": "basic",
            "contact": "@basic"
        }
    }
}
```

### *Sync* Response
#### `Searching start`
```json
{
    "operation": 0,
    "body": {}
}
```

### *Async* Responses
#### `Searching game found`
```json
{
    "operation": 2,
    "body": {
        "foundGameData": {
            "chattingStageDuration": 300000,
            "chattingTopic": "Board games",
            "choosingStageDuration": 30000,
            "localProfileId": 1,
            "profilePublicList": [
                {
                    "id": 0,
                    "username": "Major"
                },
                {
                    "id": 1,
                    "username": "Minor"
                }
            ],
            "startSessionTime": 1690309190
        }
    }
}
```

#### `Chatting new message`
```json
{
    "operation": 3,
    "body": {
        "message": {
            "senderId": 1,
            "text": "123"
        }
    }
}
```

## OPS: `Searching stop`

### Request
```json
{
    "operation": 1
}
```

## OPS: `Chatting new message`

### Request
#### `Chatting new message`
```json
{
    "operation": 3,
    "body": {
        "message": {
            "text": "Hello!"
        }
    }
}
```

### *Sync* Response
#### `Searching start`
```json
{
    "operation": 0,
    "body": {}
}
```

## OPS: `Chatting stage is over`

### *Async* Response
#### `Chatting stage is over`
```json
{
    "operation": 4,
    "body": {}
}
```

## OPS: `Choosing users chosen`

### Request
#### `Choosing users chosen`
```json
{
  "operation": 5,
  "body":{
    "userIdList": [
        1,
        2,
        3
    ]
  }
}
```

## OPS: `Choosing users chosen`

### *Async* Response
#### `Choosing users chosen`
```json
{
  "operation": 6,
  "body": {
    "matchedUsers": [
      {
        "id": 1,
        "contact": "@major"
      }
    ]
  }
}
```