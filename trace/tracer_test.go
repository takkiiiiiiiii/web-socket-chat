package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) { //*testing.T型の引数を受け取る関数はすべてユニットテストとみなされる
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値はnil")
	} else {
		tracer.Trace("こんにちは")
		if buf.String() != "こんにちは\n" {
			t.Errorf("'%s'という誤った文字列が出力", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}
