package vanilla_test

import (
	"testing"

	"github.com/dwood15/bosskeys/bk"
)

func TestLoadBasePools(t *testing.T) {
	bk.LoadBasePools("../../base_pools/oot/")
}