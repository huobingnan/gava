package classpath

import (
	"gava/jvm"
	"testing"
)

func TestCompositeClassEntry(ctx *testing.T) {
	var _ = jvm.NewClassEntry("/Users/huobingnan/env/zulu-11.jdk/Contents/Home/lib/jrt-fs.jar:/usr/local:/usr/bin")

}
