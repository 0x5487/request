# Request

A human readable and easy to use request package for Go

## Features

* human readable and easy to use
* goroutine safe
* timeout support (default is 30 seconds)
* allow to set proxy


## Usages

get request
```golang
resp, err := request.
    GET("/v1/hello").
    End()
if err != nil {
    return nil, err
}

 
if resp.OK {
    // if response's http status is between 200 ~ 299, then `resp.OK` should be ok
    // when success, do something...
    ....
}
```

get request with header

```golang
resp, err := request.
    GET("/v1/hello").
    Set("Authorization", "token") // your token
    End()
    
if err != nil {
    return nil, err
}

if resp.OK {
    // do something...
    ....
}
```

post request

```golang
form := url.Values{}
form.Add("grant_type", "password")
form.Add("client_id", "aaa")
form.Add("username", "bbb")
form.Add("password", "ccc")
body := form.Encode()

resp, err := request.
    POST("/v1/hello").
    Send(body).
    End()
if err != nil {
    return nil, err
}
if resp.OK {
    // do something...
    ....
}
```

use proxy

```golang
resp, err := request.
    GET("/v1/hello").
    SetProxyURL("http://10.2.3.4:8080").  // use proxy here
    End()
if err != nil {
    return nil, err
}
if resp.OK {
    // do something...
    ....
}
```



inspire by superAgent