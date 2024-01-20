package mysql_test

import (
	"context"
	"testing"
	"time"
	"wechat_robot/mysql"

	"github.com/stretchr/testify/assert"
)

func TestSelectHistoryMsg(t *testing.T) {
	result, err := mysql.SelectHistoryMsg(context.Background(), "@@b403c1b77f667e3382997d8471273e444dc572564253a790f9de70b112a7da0b", 10, -time.Hour)
	assert.Nil(t, err)
	t.Log(result)
}
