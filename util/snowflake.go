package util

import "github.com/sony/sonyflake"

var Snowflake = sonyflake.NewSonyflake(sonyflake.Settings{})
