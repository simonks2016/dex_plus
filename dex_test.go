package DexPlus

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/simonks2016/dex_plus/okx/param"
	"github.com/simonks2016/dex_plus/okx/rest"
)

func NewLogger() *log.Logger {
	return log.New(
		os.Stdout, // 也可以换成你自己的 Writer（SLS、文件等）
		"[OKX] ",  // 前缀
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile,
	)
}

func TestNew(t *testing.T) {

	cli := rest.NewOKXRestClient(
		rest.WithAuth(
			"",
			"",
			""),
		rest.WithSandbox(true))

	err := cli.PlaceOrder(param.PlaceOrderParams{
		InstIdCode: NewInt(3),
		TdMode:     "cash",
		ClOrdId:    NewString("a"),
		Side:       "buy",
		OrdType:    "limit",
		SZ:         "0.002",
		Px:         NewString("70000"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(5 * time.Minute)

	err = cli.CancelOrder(
		param.CancelOrder{
			InstId:  NewString("BTC-USDT"),
			ClOrdId: NewString("a"),
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func NewInt(i int) *int {
	return &i
}

func NewString(s string) *string {
	return &s
}

// 存在问题：
// 第二: 验证失败会出现批量发送验证信息，不停的重新启动
