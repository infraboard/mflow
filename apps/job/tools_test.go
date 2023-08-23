package job_test

import (
	"testing"

	"github.com/infraboard/mflow/apps/job"
)

func TestNewMapWithKVPaire(t *testing.T) {
	m := job.NewMapWithKVPaire("k1", "v1", "k2", "v2")
	t.Log(m)
}
