# WebSocket - App 🗑

## Inside the server

> Rooms have states: SEARCHING, CHATTING and CHOOSING.

## Operation (OP)

``` C++
enum class operation : int {
    SEARCHING_START        = 0, // |
    SEARCHING_STOP         = 1, // | ---> SEARCHING
    SEARCHING_GAME_FOUND   = 2, // |

    CHATTING_NEW_MESSAGE   = 3, // | ---> CHATTING
    CHATTING_STAGE_IS_OVER = 4, // |

    CHOOSING_USERS_CHOSEN  = 5, // | ---> CHOOSING
    CHOOSING_STAGE_IS_OVER = 6, // |
}
```

## SEARCHING
## OPS: `Searching start`, `Searching game found`, `Chatting new message`

### Request
#### `Searching start`
```json
{
  "operation": 0,
  "body": {
    "profile": {
      "username": "<string>",
      "contact": "<string>",
      "language": "<int>"
    }
  }
}
```

> Places the user in an available room. 
> If necessary, the server will create a room.

### *Sync* Response
#### `Searching start`
```json
{
  "operation": 0,
  "body": {}
}
```

### ... over time (see [config.yml](../../../config/config.yml)) ...

### *Async* Responses
#### `Searching game found`
```json
{
  "operation": 2,
  "body": {
    "foundGameData": {
      "chattingStageDuration": "<int>",
      "chattingTopic": "<string>",
      "choosingStageDuration": "<int>",
      "localProfileId": "<int>",
      "profilePublicList": [
        {
          "id": "<int>",
          "username": "<string>"
        },
        ...
      ],
      "startSessionTime": "<int>"
    }
  }
}
```

> Users are assigned a local number 
> in the order they are added to the room.

## OPS: `Searching stop`

### Request
```json
{
  "operation": 1
}
```

> Until information about the found game comes.
> Or just gracefully close the websocket.

## CHATTING
## OPS: `Chatting new message`, `Chatting stage is over`

### Request
#### `Chatting new message`
```json
{
  "operation": 3,
  "body": {
    "message": {
      "text": "<string>"
    }
  }
}
```

### *Sync* or *Async* Response
#### `Chatting new message`
```json
{
  "operation": 3,
  "body": {
    "message": {
      "senderId": "<int>",
      "text": "<string>"
    }
  }
}
```

> Your own message will arrive almost in sync.

### ... over time (see [config.yml](../../../config/config.yml)) ...

### *Async* Response
#### `Chatting stage is over`
```json
{
  "operation": 4,
  "body": {}
}
```

## CHOOSING
## OPS: `Choosing users chosen`

### Request
#### `Choosing users chosen`
```json
{
  "operation": 5,
  "body": {
    "userIdList": [
      "<int>",
      ...
    ]
  }
}
```

### ... over time (see [config.yml](../../../config/config.yml)) ...

### *Async* Response
#### `Choosing users chosen`
```json
{
  "operation": 6,
  "body": {
    "matchedUsers": [
      {
        "id": "<int>",
        "contact": "<string>"
      },
      ...
    ]
  }
}
```