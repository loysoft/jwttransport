Go JWT HTTP Transport
=====================

Go HTTP transport with JWT (OAUth2) Authorization support.


Installation
------------

```shell
go get github.com/loysoft/jwttransport
go install github.com/loysoft/jwttransport
```

Usage
-----
TODO

Usage examples:

config-file-with-secrets.json:
```json
{
  "private_key_id": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIC...-----END PRIVATE KEY-----\n",
  "client_email": "XXXXXXXXXXXX-YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY@developer.gserviceaccount.com",
  "client_id": "XXXXXXXXXXXX-YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY.apps.googleusercontent.com",
  "type": "service_account"
}
```

main.go:
```go
import (
  jwtt "github.com/loysoft/jwttransport"
  pubsub "code.google.com/p/google-api-go-client/pubsub/v1beta1"
)

...

var err error

co := &jwtt.Configurator{
	Scope: "https://www.googleapis.com/auth/pubsub https://www.googleapis.com/auth/cloud-platform",
}

jwttConfig, err = co.Load("config-file-with-secrets.json")
if err != nil {
	return err
}

transport := &jwtt.Transport{Config: jwttConfig}

err = transport.PrepareToken()
if err != nil {
	return err
}

client := transport.Client()

pubsubService, err := pubsub.New(client)
if err != nil {
	return err
}

pubsubService.Topics.List().Query(...
```
