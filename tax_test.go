package tax

import "testing"

func TestGuildStartsAtZero(t *testing.T) {
	if New(30, 1000).Guild() != 0 {
		t.Fatal("公会税收初始应该是 0")
	}
}

func TestEarnKeepsConservation(t *testing.T) {
	ledger := New(30, 1000)
	charged, net := ledger.Earn("p1", 100, "r1")
	if charged+net != 100 {
		t.Fatalf("税加实得应该等于收入：%v %v", charged, net)
	}
}

func TestTodayStartsAtZero(t *testing.T) {
	if New(30, 1000).Today("p1") != 0 {
		t.Fatal("今日已缴初始应该是 0")
	}
}

func TestRateIsStored(t *testing.T) {
	if New(25, 1000).ratePct != 25 {
		t.Fatal("比例应该能读回来")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New(30, 1000).Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestEarnRoundsDown(t *testing.T) {
	ledger := New(30, 1000)
	charged, net := ledger.Earn("p1", 150, "r1")
	if charged != 45 {
		t.Fatalf("150 按 30%% 向下取整应该是 45，得到 %d", charged)
	}
	if net != 105 {
		t.Fatalf("实得应该是 105，得到 %d", net)
	}
	if charged+net != 150 {
		t.Fatal("税加实得应该恒等于收入")
	}
}

func TestDailyCapStopsCharging(t *testing.T) {
	ledger := New(30, 10)
	first, _ := ledger.Earn("p1", 100, "r1")
	second, net := ledger.Earn("p1", 100, "r2")
	if first != 10 {
		t.Fatalf("首笔税应该被上限截到 10，得到 %d", first)
	}
	if second != 0 {
		t.Fatalf("超出每日上限后不应再收税，得到 %d", second)
	}
	if net != 100 {
		t.Fatalf("超出上限部分实得应相应变多，得到 %d", net)
	}
	if ledger.Today("p1") != 10 {
		t.Fatal("当日累计缴税不应超过 dailyCap")
	}
	if ledger.Guild() != 10 {
		t.Fatal("公会税收只应包含实际收到的 10")
	}
}

func TestExemptMemberPaysNoTax(t *testing.T) {
	ledger := New(30, 1000)
	ledger.Exempt("p1")
	charged, net := ledger.Earn("p1", 100, "r1")
	if charged != 0 {
		t.Fatalf("豁免成员不应被收税，得到 %d", charged)
	}
	if net != 100 {
		t.Fatalf("豁免成员实得应为全额，得到 %d", net)
	}
	if ledger.Guild() != 0 {
		t.Fatal("豁免成员不应增加公会税收")
	}
}

func TestDuplicateReqIDIsIdempotent(t *testing.T) {
	ledger := New(30, 1000)
	firstTax, firstNet := ledger.Earn("p1", 100, "same")
	secondTax, secondNet := ledger.Earn("p1", 100, "same")
	if secondTax != firstTax || secondNet != firstNet {
		t.Fatalf("重复 reqID 应返回上次结果：(%d,%d) vs (%d,%d)",
			firstTax, firstNet, secondTax, secondNet)
	}
	if ledger.Guild() != 30 {
		t.Fatalf("重复请求不应改变公会税收，得到 %d", ledger.Guild())
	}
	if ledger.Today("p1") != 30 {
		t.Fatal("重复请求不应重复累计当日税额")
	}
}

func TestScannedDoesNotGrowWithMembers(t *testing.T) {
	ledger := New(30, 1000000)
	for index := 0; index < 100000; index++ {
		ledger.Earn("p", 100, "r0")
	}
	baseline := ledger.Scanned()
	for index := 0; index < 100000; index++ {
		ledger.Earn("p", 100, "r0")
	}
	if ledger.Scanned() != baseline+100000 {
		t.Fatalf("Scanned 每次结算应摊还 O(1)，baseline=%d now=%d",
			baseline, ledger.Scanned())
	}
	if ledger.Guild() != 30 {
		t.Fatal("同一 reqID 结算十万次仍应只收一次税")
	}
}
