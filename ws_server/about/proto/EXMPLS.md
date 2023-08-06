# Packet Examples

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

### *Async* Response
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

### Request
#### `Searching stop`
```json
{
  "operation": 1
}
```

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

### *Async* Response
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

### *Async* Response
#### `Chatting stage is over`
```json
{
  "operation": 4,
  "body": {}
}
```

### Request
#### `Choosing users chosen`
```json
{
  "operation": 5,
  "body": {
    "userIdList": [
      1,
      2,
      3
    ]
  }
}
```

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