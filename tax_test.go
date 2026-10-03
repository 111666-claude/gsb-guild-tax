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
