package clone

var (
	DefaultConfig = &Config{
		BufSize:   100,
		Threshold: 100,
	}
)

type Config struct {
	// AddFile時のチャネルのバッファーサイズ(デフォルト:100)
	BufSize int
	// 連続トークン数の境界値(デフォルト:100)
	Threshold int
	Serializer
	SuffixTree
}
