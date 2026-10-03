# guild-tax

公会税收：按比例抽成、每日上限、豁免名单与请求幂等，只用 Go 标准库。

```
go run ./cmd/guild-tax --sample=round
go run ./cmd/guild-tax --sample=cap
go run ./cmd/guild-tax --sample=exempt
go run ./cmd/guild-tax --sample=idem
go run ./cmd/guild-tax --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **抽成**：税额等于 `amount * ratePct / 100`，向下取整。
- **守恒**：任何一次收税都满足 `税 + 实得 == amount`。
- **每日上限**：同一个成员一天内累计缴纳的税不超过 `dailyCap`，超出部分不再收，
  实得相应变多。
- **豁免名单**：名单里的成员一律不收税。
- **幂等**：同一个 `reqID` 只结算一次，重复调用返回上次结果且不改变公会税收。
- **不变量**：公会税收只增不减；同一成员当日已缴不超过上限；
  同一串收入序列重复执行结果相同。
- **代价**：收税与查询不许遍历全部成员，`Scanned` 不随成员数线性增长。
- **规模**：10 万成员、每秒 10 万次结算，单次摊还 O(1)。

## 输出契约（不改格式）

```
tax=..
second=..
guild=..
scanned=..
```

## 目录

```
tax.go              抽成、上限、豁免与幂等
tax_test.go         单元测试
cmd/guild-tax/      命令行入口
```
