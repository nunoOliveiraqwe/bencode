package bencode

import (
	"testing"
)

func TestEncode(t *testing.T) {
	dict := make(map[string]interface{})
	dict["announce"] = "https://torrent.ubuntu.com/announce"
	dict["announce-list"] = []interface{}{
		[]interface{}{
			"https://torrent.ubuntu.com/announce",
		},
		[]interface{}{
			"https://ipv6.torrent.ubuntu.com/announce",
		},
	}
	dict["comment"] = "Ubuntu CD releases.ubuntu.com"
	dict["created by"] = "mktorrent 1.1"
	dict["creation date"] = 1634219565

	info := make(map[string]interface{})
	info["length"] = 3116482560
	info["name"] = "ubuntu-21.10-desktop-amd64.iso"
	info["piece length"] = 262144
	info["pieces"] = "earls"
	dict["info"] = info

	res := string(Encode(dict))
	expected := "d8:announce35:https://torrent.ubuntu.com/announce13:announce-listll35:https://torrent.ubuntu.com/announceel40:https://ipv6.torrent.ubuntu.com/announceee7:comment29:Ubuntu CD releases.ubuntu.com10:created by13:mktorrent 1.113:creation datei1634219565e4:infod6:lengthi3116482560e4:name30:ubuntu-21.10-desktop-amd64.iso12:piece lengthi262144e6:pieces5:earlsee"
	if res != expected {
		t.Errorf("expected %s\ngot %s", expected, res)
	}
}
