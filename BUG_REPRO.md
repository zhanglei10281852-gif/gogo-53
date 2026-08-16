# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

某个组件因为健康指标不达标 apply 失败，修好指标后重跑 apply，状态里该组件的迁移记录变成了同一个 ID 出现两次，重试次数越多重复越多。请修复 apply 的迁移记录逻辑，让重试不产生重复条目，同时保持首次 apply 的迁移集合和顺序不变，并保证全量测试通过。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-53
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-53.git
- parent SHA：ae38916b3605c16665314d4b6f454238d064e125

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-53.git bug-repro
cd bug-repro
git checkout --detach ae38916b3605c16665314d4b6f454238d064e125
go test ./internal/rail -run "^TestApplySimulationDoesNotDuplicateMigrationsOnRetry$" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/rail -run "^TestApplySimulationDoesNotDuplicateMigrationsOnRetry$" -count=1 -v
=== RUN   TestApplySimulationDoesNotDuplicateMigrationsOnRetry
    reapply_regression_test.go:36: retried apply duplicated migration records: [schema schema]
--- FAIL: TestApplySimulationDoesNotDuplicateMigrationsOnRetry (0.00s)
FAIL
FAIL	releaserail/internal/rail	0.002s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/rail -run "^TestApplySimulationDoesNotDuplicateMigrationsOnRetry$" -count=1 -v
=== RUN   TestApplySimulationDoesNotDuplicateMigrationsOnRetry
    reapply_regression_test.go:36: retried apply duplicated migration records: [schema schema]
--- FAIL: TestApplySimulationDoesNotDuplicateMigrationsOnRetry (0.04s)
FAIL
FAIL	releaserail/internal/rail	0.160s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

失败后重试 apply 得到的迁移记录不含重复 ID，且内容与单次成功 apply 一致；波次、状态流转与健康判定不回归；双架构定向、全量、build/vet 通过。
