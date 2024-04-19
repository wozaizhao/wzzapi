package global

import (
	"testing"
)

func TestAbs(t *testing.T) {
	fileName := GetFileNameFromUrl("https://d.wozaizhao.com/logo.png")
	if fileName != "logo.png" {
		t.Errorf("fileName is %s", fileName)
	}
	f2 := GetFileNameFromUrl("https://img.alicdn.com/imgextra/i2/O1CN01FF1t1g1Q3PDWpSm4b_!!6000000001920-55-tps-508-135.svg")
	if f2 != "O1CN01FF1t1g1Q3PDWpSm4b_!!6000000001920-55-tps-508-135.svg" {
		t.Errorf("fileName is %s", fileName)
	}
	f3 := GetFileNameFromUrl("https://d.wozaizhao.com/logo.png?r21r12r=3211234321")
	if f3 != "logo.png" {
		t.Errorf("fileName is %s", fileName)
	}
}
