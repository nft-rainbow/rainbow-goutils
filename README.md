# rainbow-goutils

## Alert

### Config Sample
```yaml
alert:
  customTags: [alert, local]
  channels:
    default:
      platform: dingtalk      
      webhook: "YOUR_WEBHOOK"
      secret: ""
      atMobiles: []
      isAtAll: false
    channel2:
      platform: dingtalk
      webhook: "YOUR_WEBHOOK_2"
      secret: ""
      atMobiles: ["13111112222"]
      isAtAll: false
```

### Usage

```go
alertutils.MustInitFromViper()
alertutils.DingInfof("test")
alertutils.DingWarnf("test %s", "test")
```
