// Package tax 是公会税收：按比例抽成、每日上限、豁免名单与请求幂等。
// 缺陷：抽成用四舍五入、每日上限不生效、豁免名单不生效、请求不幂等、
// 每次收税遍历全部成员。
package tax

type record struct {
	tax int
	net int
}

// Ledger 是公会税收台账。
type Ledger struct {
	ratePct  int
	dailyCap int
	exempt   map[string]bool
	today    map[string]int
	guild    int
	memo     map[string]record
	scanned  int
}

// New 建台账。
func New(ratePct, dailyCap int) *Ledger {
	return &Ledger{ratePct: ratePct, dailyCap: dailyCap, exempt: map[string]bool{},
		today: map[string]int{}, memo: map[string]record{}}
}

// Exempt 把成员加入免税名单。
func (l *Ledger) Exempt(member string) { l.exempt[member] = true }

// Earn 记一次成员收入，返回 (税, 实得)。
func (l *Ledger) Earn(member string, amount int, reqID string) (int, int) {
	l.scanned += len(l.today)
	tax := (amount + 50) / 100 * l.ratePct
	l.guild += tax
	l.today[member] += tax
	l.memo[reqID] = record{tax: tax, net: amount - tax}
	return tax, amount - tax
}

// Guild 是公会累计税收。
func (l *Ledger) Guild() int { return l.guild }

// Today 是某成员今天已缴的税。
func (l *Ledger) Today(member string) int { return l.today[member] }

// Scanned 是累计扫描量（规模观测）。
func (l *Ledger) Scanned() int { return l.scanned }
