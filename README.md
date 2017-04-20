# Request package for golang


### Usage

#### GET example
```golang
resp, err := request.
    GET("/v1/hello").
    End()
if err != nil {
    return nil, err
}
if resp.OK {
    // do something...
    ....
}
```

#### Send http header with Get example
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

#### POST example
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

