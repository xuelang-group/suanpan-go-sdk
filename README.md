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
storage.DeleteMultiObjects(objectNames)
```

### 4. Log Api

```go
import "github.com/xuelang-group/suanpan-go-sdk/suanpan/log"

log.Trace("trace message")
log.Tracef("trace message: %s", msg)
log.Debug("debug message")
log.Debugf("debug message: %s", msg)
log.Info("info message")
log.Infof("info message: %s", msg)
log.Warn("warn message")
log.Warnf("warn message: %s", msg)
log.Error("error message")
log.Errorf("error message: %s", msg)
```
