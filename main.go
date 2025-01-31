package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

func main() {
	router := gin.Default()

	// 启用CORS
	router.Use(corsMiddleware())

	// 静态文件服务
	router.Static("/public", "./public")
	
	// 启用CORS
	router.Use(corsMiddleware())

	// API端点：获取public/2025md目录下的指定Markdown文件
	router.GET("/api/markdown/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join("public", "2025md", filename)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		c.Data(http.StatusOK, "text/plain", data)
	})

	// API端点：获取public/2025pic目录下的指定图片文件
	router.GET("/api/image/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join("public", "2025pic", filename)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		c.Data(http.StatusOK, "image/"+filepath.Ext(filename)[1:], data)
	})

	// API端点：获取歌曲信息
	router.GET("/api/songInfo", func(c *gin.Context) {
		songId := c.Query("songId")
		songUrl := fmt.Sprintf("https://music.163.com/api/v3/song/detail?c=[{\"id\":%s}]", songId)

		body, err := fetch(songUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Data(http.StatusOK, "application/json", body)
	})

	// API端点：获取歌词
	router.GET("/api/lyricsInfo", func(c *gin.Context) {
		songId := c.Query("songId")
		lyricsUrl := fmt.Sprintf("https://music.163.com/api/song/media?id=%s", songId)

		body, err := fetch(lyricsUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Data(http.StatusOK, "application/json", body)
	})

	// API端点：获取专辑封面
	router.GET("/api/picInfo", func(c *gin.Context) {
		url := c.Query("picUrl")

		body, err := fetch(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Data(http.StatusOK, "image/jpeg", body)
	})

	// API端点：获取歌单
	router.GET("/api/playlist", func(c *gin.Context) {
		filePath := filepath.Join("public", "2024song", "playlist.json")
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	})

	// API端点：获取flac
// API端点：获取flac或mp3
router.GET("/api/flac", func(c *gin.Context) {
    songId := c.Query("songId")
    jsonPath := filepath.Join("public", "2024song", "playlist.json")

    data, err := ioutil.ReadFile(jsonPath)
    if err != nil {
        c.String(http.StatusNotFound, "File not found")
        return
    }

    var songs []map[string]interface{}
    if err := json.Unmarshal(data, &songs); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    for _, song := range songs {
        if fmt.Sprint(song["id"]) == songId {
            songName, ok := song["name"].(string)
            if !ok {
                c.String(http.StatusInternalServerError, "Error processing song data")
                return
            }
            flacPath := filepath.Join("public", "2024song", "music", songName+".flac")
            mp3Path := filepath.Join("public", "2024song", "music", songName+".mp3")

            if _, err := os.Stat(flacPath); err == nil {
                c.File(flacPath)
                return
            } else if _, err := os.Stat(mp3Path); err == nil {
                c.File(mp3Path)
                return
            }
        }
    }

    // If no valid song is found
    c.String(http.StatusNotFound, "Song not found")
})


	// 启动服务器
	router.Run(":4000")
}

func fetch(url string) ([]byte, error) {
	statusCode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status: %d", statusCode)
	}
	return body, nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Next()
	}
}