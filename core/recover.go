package core

import sctx "github.com/ngothiep2412/rockship-go-sdk"

func Recover() {
	if r := recover(); r != nil {
		sctx.GlobalLogger().GetLogger("recovered").Errorln(r)
	}
}
