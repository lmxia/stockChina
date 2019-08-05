package types
//StockPriceInfo 股票价格信息
type StockPriceHistoryInfo struct {
	Date 				  string  // 日期：yymmdd
	OpeningPriceToday     float64 // 今日开盘价
	ClosingPriceYesterday float64 // 昨日收盘价
	HighestPriceToday     float64 // 今日最高价
	LowestPriceToday      float64 // 今日最低价
}
