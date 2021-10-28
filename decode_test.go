package bencode

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	str := "d8:announce35:https://torrent.ubuntu.com/announce13:announce-listll35:https://torrent.ubuntu.com/announceel40:https://ipv6.torrent.ubuntu.com/announceee7:comment29:Ubuntu CD releases.ubuntu.com10:created by13:mktorrent 1.113:creation datei1634219565e4:infod6:lengthi3116482560e4:name30:ubuntu-21.10-desktop-amd64.iso12:piece lengthi262144e6:pieces5:earlsee"
	dict, err := Decode(strings.NewReader(str))
	if err != nil {
		t.Error(err)
	}
	if dict["announce"] != "https://torrent.ubuntu.com/announce" {
		t.Error("announce mismatch")
	} else if dict["comment"] != "Ubuntu CD releases.ubuntu.com" {
		t.Error("comment mismatch")
	} else if len(dict["announce-list"].([]interface{})) != 2 {
		t.Error("invalid announce list size")
	} else if dict["creation date"].(int64) != 1634219565 {
		t.Error("creation date mismatch")
	}
}

func TestDecodeNegativeLenght(t *testing.T) {
	_, err := Decode(strings.NewReader("d3:key-1:e"))
	if err == nil {
		t.Error("string length can not be a negative number")
	}
}

func TestDecodeZeroLength(t *testing.T) {
	_, err := Decode(strings.NewReader("d3:key0:e"))
	if err != nil {
		t.Error("error while non expected")
	}

}
