package proxy

import (
	"EZ-Encrypt-Middleware/config"
	"EZ-Encrypt-Middleware/utils"
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ProxyHandler handles incoming requests and forwards them to the backend API
func ProxyHandler(c *gin.Context) {
	fullPath := c.Request.URL.Path

	encodedPath := strings.TrimPrefix(fullPath, "/")

	iv := c.GetHeader("X-IV")
	if iv == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少X-IV请求头"})
		return
	}

	encryptedPath, err := base64.URLEncoding.DecodeString(encodedPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Base64编码路径"})
		return
	}

	decryptedPath := utils.Decryption(string(encryptedPath), iv)

	backendURL := config.AppConfig.BackendAPIURL

	targetURL := backendURL + "/api/v1" + decryptedPath

	target, err := url.Parse(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的目标URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	timeout := 30 * time.Second
	if config.AppConfig.RequestTimeout != "" {
		if t, err := time.ParseDuration(config.AppConfig.RequestTimeout + "ms"); err == nil {
			timeout = t
		}
	}

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// 修改请求以转发到解密后的路径
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.URL.RawQuery = target.RawQuery
		req.Host = target.Host

		req.Header.Del("X-IV")

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		req = req.WithContext(ctx)
		defer cancel()
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Credentials")
		resp.Header.Del("Access-Control-Allow-Headers")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Max-Age")

		return nil
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
