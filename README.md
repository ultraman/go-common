# go-common

```bazaar
go get github.com/yaoliu/go-common@v0.1-beta
```

## rest 主要借鉴/抄袭/复制于client-go rest实现

```bazaar
import (
    "github.com/yaoliu/go-common/rest"
)

func main(){
    config := &rest.Config{
		Host:      "https://localhost",
		UserAgent: "go-user-agent",
		ContentConfig: rest.ContentConfig{
			Codec:       rest.NewJsonMarshaler(),
			ContentType: "application/json;charset=UTF-8",
			Timeout:     10,
		},
	}
	
    client, err := rest.NewRESTClientFor(config)
    if err != nil {
        return nil
    }
    client.Get().Path("/api/v1/user/info").Timeout(10).Do(ctx).Raw()
}
```
