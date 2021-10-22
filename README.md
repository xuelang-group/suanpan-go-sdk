# suanpan-go-sdk

## 功能列表

### 1. 获得右面板配置参数

```go
import "github.com/xuelang-group/suanpan-go-sdk/suanpan/parameter"

parameter.Get("--xxx")
```

### 2. 发送消息与接受消息

```go
import "github.com/xuelang-group/suanpan-go-sdk/suanpan/stream"

func handle(r stream.Request) {
        r.Send(map[string]interface{}{
                "out1": r.Data,
        })
}
```

### 3. Storage Api

```go
import "github.com/xuelang-group/suanpan-go-sdk/suanpan/storage"

storage.FGetObject(objectName, filePath)
storage.FPutObject(objectName, filePath)
storage.PutObject(objectName, data)
storage.ListObjects(prefix)
storage.DeleteObject(objectName)
storage.DeleteObjects(objectNames)
```
