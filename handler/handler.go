package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	types "stockChina/type"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StockWatchInfo struct {
	StockID string `json:"stock_id"`
}

type BaseURLResolver interface {
	Resolve(c *gin.Context) string
}

type TensentStockURLResolver struct {
	ExternalUrl string
}

// NewHTTPClientReverseProxy proxies to an upstream host through the use of a http.Client
func NewHTTPClientReverseProxy(timeout time.Duration, maxIdleConns, maxIdleConnsPerHost int) *HTTPClientReverseProxy {
	h := HTTPClientReverseProxy{
		Timeout: timeout,
	}

	h.Client = http.DefaultClient
	h.Timeout = timeout
	h.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	h.Client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          maxIdleConns,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &h
}

// HTTPClientReverseProxy proxy to a remote BaseURL using a http.Client
type HTTPClientReverseProxy struct {
	TargetUrl string
	Client    *http.Client
	Timeout   time.Duration
}

// Resolve used tensent resolver
func (resolver TensentStockURLResolver) Resolve(c *gin.Context) string {
	year, _ := strconv.ParseInt(c.Query("year"), 10, 32)
	address := c.Query("address")
	stockId := c.Query("stockid")
	return fmt.Sprintf("%s%d/%s%s.js", resolver.ExternalUrl, year, address, stockId)
}

func MakeHandlerWrapper(next gin.HandlerFunc, resovler BaseURLResolver, proxy *HTTPClientReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stockInfos []types.StockPriceHistoryInfo
		baseURL := resovler.Resolve(c)
		proxy.TargetUrl = baseURL
		res, resErr := forwardRequest(c, proxy)
		_all, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		_allString := string(_all)
		allInfo := strings.Split(_allString, "\n")
		for i, dayInfo := range allInfo {
			if i == 0 {
				continue
			}
			infos := strings.Split(dayInfo, " ")
			if len(infos) != 6 {
				break
			}
			openPrice, _ := strconv.ParseFloat(infos[1], 64)
			closePrice, _ := strconv.ParseFloat(infos[2], 64)
			maxPrice, _ := strconv.ParseFloat(infos[3], 64)
			minPrice, _ := strconv.ParseFloat(infos[4], 64)
			stockPriceHistoryInfo := types.StockPriceHistoryInfo{
				Date:                  infos[0],
				OpeningPriceToday:     openPrice,
				ClosingPriceYesterday: closePrice,
				HighestPriceToday:     maxPrice,
				LowestPriceToday:      minPrice,
			}
			stockInfos = append(stockInfos, stockPriceHistoryInfo)
		}

		if resErr != nil {
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": stockInfos})
		next(c)
	}
}

func forwardRequest(c *gin.Context, proxy *HTTPClientReverseProxy) (*http.Response, error) {
	proxyClient := proxy.Client
	upstreamReq, err := http.NewRequest(c.Request.Method, proxy.TargetUrl, nil)
	if err != nil {
		return nil, nil
	}
	deleteHeaders(&upstreamReq.Header, &hopHeaders)
	if len(c.Request.Host) > 0 && upstreamReq.Header.Get("X-Forwarded-Host") == "" {
		upstreamReq.Header["X-Forwarded-Host"] = []string{c.Request.Host}
	}
	if upstreamReq.Header.Get("X-Forwarded-For") == "" {
		upstreamReq.Header["X-Forwarded-For"] = []string{c.Request.RemoteAddr}
	}

	if c.Request.Body != nil {
		upstreamReq.Body = c.Request.Body
		defer upstreamReq.Body.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), proxy.Timeout)
	defer cancel()

	return proxyClient.Do(upstreamReq.WithContext(ctx))
}

func copyHeaders(destination http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(destination)[k] = vClone
	}
}

func deleteHeaders(target *http.Header, exclude *[]string) {
	for _, h := range *exclude {
		target.Del(h)
	}
}

// AggregateAudits handle aggres handler
func GetKline(c *gin.Context) {
	return
}

var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}
