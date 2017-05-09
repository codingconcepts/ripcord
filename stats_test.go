package ripcord

import (
	"testing"

	"github.com/codingconcepts/albert/test"
)

func TestFilter(t *testing.T) {
	in := IOStats{
		IOStat{Name: "a"},
		IOStat{Name: "b"},
		IOStat{Name: "c"},
		IOStat{Name: "d"},
		IOStat{Name: "e"},
	}

	out := in.Filter("b", "d")

	test.Equals(t, 2, len(out))
	test.Equals(t, "b", out[0].Name)
	test.Equals(t, "d", out[1].Name)
}

func TestFind(t *testing.T) {
	in := IOStats{
		IOStat{Name: "a"},
		IOStat{Name: "b"},
		IOStat{Name: "c"},
		IOStat{Name: "d"},
		IOStat{Name: "e"},
	}

	out := in.Find("c")

	test.Equals(t, "c", out.Name)
}

func TestFindNothing(t *testing.T) {
	in := IOStats{
		IOStat{Name: "a"},
	}

	out := in.Find("b")

	test.Equals(t, "", out.Name)
}
