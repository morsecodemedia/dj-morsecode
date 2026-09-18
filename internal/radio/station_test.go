package radio

import "testing"

func TestFindByStreamURL(t *testing.T) {

	station, ok := FindByStreamURL(
		"https://wxpn.xpn.org/xpnmp3hi",
	)
	if !ok {
		t.Fatal("expected station to be found")
	}

	if station.ID != "wxpn" {
		t.Errorf(
			"expected station ID %q, got %q",
			"wxpn",
			station.ID,
		)
	}

}

func TestFindByStreamURLMiss(t *testing.T) {

	_, ok := FindByStreamURL(
		"https://example.com/not-a-station",
	)

	if ok {
		t.Fatal("expected station lookup to miss")
	}

}
