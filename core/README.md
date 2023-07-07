# 公用库

## 使用说明
### mq 新增 topic
1. 在 core/config/config.go 文件中，新增一个返回 topic 的函数，如：

```go
// 同步飞书用户、部门信息 topic
func GetMqSyncMemberDeptTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.SyncMemberDept
}
```

2.在 core/config/config.go 中的 TopicConfig 结构体中新增一个你需要的 topic 配置，如 `SyncMemberDept TopicConfigInfo`。

3.在实际的业务项目的配置文件中，配置对应的值，以 polaris-backend 为例，在 config/application.common.yaml 中的 `MQ.Topics` 下配置 SyncMemberDept 对应的值：

```
SyncMemberDept:
  Topic: topic_org_sync_member_dept_local
  GroupId: topic_org_sync_member_dept_group_local
```

3.提交打包，发布新版本。
