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

func TestTaxFloorsToInteger(t *testing.T) {
	ledger := New(30, 1000)
	charged, net := ledger.Earn("p1", 150, "r1")
	if charged != 45 {
		t.Fatalf("150*30%% 整除向下取整应为 45，实际 %d", charged)
	}
	if charged+net != 150 {
		t.Fatalf("税加实得应恒等于收入：%d+%d", charged, net)
	}

	ledger2 := New(30, 1000)
	charged2, _ := ledger2.Earn("p2", 3, "r2")
	if charged2 != 0 {
		t.Fatalf("3*30%% 向下取整应为 0，实际 %d", charged2)
	}
}

func TestDailyCapStopsCharging(t *testing.T) {
	ledger := New(30, 10)
	first, firstNet := ledger.Earn("p1", 100, "r1")
	if first != 10 || firstNet != 90 {
		t.Fatalf("首笔受每日上限应为 10/90，实际 %d/%d", first, firstNet)
	}
	second, secondNet := ledger.Earn("p1", 100, "r2")
	if second != 0 || secondNet != 100 {
		t.Fatalf("超上限后不应再收税，实得 100，实际 %d/%d", second, secondNet)
	}
	if ledger.Today("p1") != 10 {
		t.Fatalf("当日累计不应超过 dailyCap=10，实际 %d", ledger.Today("p1"))
	}
	if ledger.Guild() != 10 {
		t.Fatalf("公会累计税收应为 10，实际 %d", ledger.Guild())
	}
}

func TestDailyCapIsPerMember(t *testing.T) {
	ledger := New(30, 10)
	ledger.Earn("p1", 100, "r1")
	charged, _ := ledger.Earn("p2", 100, "r2")
	if charged != 10 {
		t.Fatalf("上限应按成员独立计算，实际 %d", charged)
	}
}

func TestExemptMemberPaysNoTax(t *testing.T) {
	ledger := New(30, 1000)
	ledger.Exempt("p1")
	charged, net := ledger.Earn("p1", 100, "r1")
	if charged != 0 || net != 100 {
		t.Fatalf("豁免成员不应被收税，实际 %d/%d", charged, net)
	}
	if ledger.Guild() != 0 || ledger.Today("p1") != 0 {
		t.Fatal("豁免收入不应计入公会税收和当日累计")
	}
}

func TestRepeatedReqIDIsIdempotent(t *testing.T) {
	ledger := New(30, 1000)
	tax1, net1 := ledger.Earn("p1", 100, "same")
	tax2, net2 := ledger.Earn("p1", 100, "same")
	if tax1 != tax2 || net1 != net2 || tax1 != 30 {
		t.Fatalf("重复 reqID 应返回上次结果 30/70，实际 %d/%d", tax2, net2)
	}
	if ledger.Guild() != 30 {
		t.Fatalf("重复结算不应改变公会税收，实际 %d", ledger.Guild())
	}
	if ledger.Today("p1") != 30 {
		t.Fatalf("重复结算不应改变当日累计，实际 %d", ledger.Today("p1"))
	}
}

func TestScannedIsConstantPerCall(t *testing.T) {
	ledger := New(30, 1000000)
	for index := 0; index < 3000; index++ {
		ledger.Earn("p-"+string(rune('a'+index%26))+string(rune('0'+index/26%10)), 100,
			"r-"+string(rune('a'+index%26))+string(rune('0'+index/26%10)))
	}
	for index := 0; index < 3000; index++ {
		ledger.Today("p-z9")
	}
	if scanned := ledger.Scanned(); scanned > 6000 {
		t.Fatalf("Scanned 不应随成员数线性增长，实际 %d", scanned)
	}
}
