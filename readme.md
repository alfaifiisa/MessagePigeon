# MessagePigeon

MessagePigeon is a project that is intended for demonstration purposes only

Resource:
- https://en.wikipedia.org/wiki/Concatenated_SMS
- http://mobiletidings.com/2009/02/18/combining-sms-messages/

#### Running
default port is 8080

``` bash
    go build
    ./MessagePigeon -port 7070
```

### API documentaion:


##### /messages
this url is for sending a new sms message

```http
POST /messages HTTP/1.1
Content-Type: application/json

{
  "recipient":232222222222,
  "originator":"33",
  "message":"This is a test messag"
}
```
