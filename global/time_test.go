package global_test

import (
	"testing"
	"time"

	"github.com/vkuznecovas/mouthful/global"

	"github.com/stretchr/testify/assert"
)

func TestTimeFromNano(t *testing.T) {
	timeNow := time.Now()
	nano := timeNow.UnixNano()
	time := global.NanoToTime(nano)
	assert.Equal(t, time.UnixNano(), timeNow.UnixNano())
}
