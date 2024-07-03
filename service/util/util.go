package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func IsWindows() bool {
	if strings.Compare(runtime.GOOS, "windows") == 0 {
		return true
	}
	return false
}

func GetFileAndExtensionFromURL(rawurl string) (string, string) {
	// 解析URL
	u, err := url.Parse(rawurl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", ""
	}

	// 获取URL的路径部分
	path := u.Path

	// 从路径中提取文件名,包含了后缀
	filename := filepath.Base(path)

	// 获取文件扩展名
	ext := filepath.Ext(filename)

	return filename, ext
}

// DownloadFile 下载文件并保存到本地
func DownloadFile(filedir string, url string) (string, error) {
	// 创建HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查HTTP请求是否成功
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned HTTP status %v", resp.StatusCode)
	}

	// 获取文件名和扩展名
	fileName, extension := GetFileAndExtensionFromURL(url)
	fmt.Println(filedir, string(filepath.Separator), fileName, extension)
	filePath := fmt.Sprintf("%s%s%s", filedir, string(filepath.Separator), fileName)
	fmt.Println(filePath)
	// 打开文件准备写入
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 将下载的数据写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
