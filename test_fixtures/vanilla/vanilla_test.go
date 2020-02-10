package vanilla_test

import (
	"testing"

	bk "github.com/dwood15/bosskeys"
)

func TestLoadBasePools(t *testing.T) {
	bk.LoadBasePools("../../base_pools/oot/")
}