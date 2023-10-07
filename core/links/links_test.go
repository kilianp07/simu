package links

import "testing"

const (
	csvPath        = "../../test/core/links/links.csv"
	source_adapter = "sender"
	source_key     = "value"
	target_adapter = "receiver"
	target_key     = "value"
)

func TestNew(t *testing.T) {
	var (
		l   *Links
		err error
	)
	if l, err = New(csvPath, nil); err != nil {
		t.Fatal(err)
	}

	if l.Links[0].source.adapterName != source_adapter {
		t.Fatal("source adapter name is not correct, expected: ", source_adapter, ", got: ", l.Links[0].source.adapterName)
	}

	if l.Links[0].source.key != source_key {
		t.Fatal("source key is not correct, expected: ", source_key, ", got: ", l.Links[0].source.key)
	}

	if l.Links[0].target.adapterName != target_adapter {
		t.Fatal("target adapter name is not correct, expected: ", target_adapter, ", got: ", l.Links[0].target.adapterName)
	}

	if l.Links[0].target.key != target_key {
		t.Fatal("target key is not correct, expected: ", target_key, ", got: ", l.Links[0].target.key)
	}

}
