package global

var (
	BlogInfo *Blog
)

func init() {
	BlogInfo = new(Blog)
}

type Blog struct {
	DefaultPageNum int
	BlogName       string
	BTitle         string
	SubTitle       string
	BeiAn          string
	CopyRight      string
}
