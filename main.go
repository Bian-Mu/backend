package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)


func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Next()
	}
}

func fetch(url string) ([]byte, error) {
	statusCode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("%d", statusCode)
	}
	return body, nil
}

func main(){
	server:=gin.Default()

	//跨域问题
	server.Use(corsMiddleware())

	publicGroup:=server.Group("/public")
	{
		publicGroup.GET("/2025md/:mdname",func(c *gin.Context){
			mdname:=c.Param("mdname")
			mdpath:=filepath.Join("public","2025md",mdname)
			data,err:=os.ReadFile(mdpath)
			if err != nil{
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Data(http.StatusOK,"text/plain",data)
		})

		publicGroup.GET("/2025pic/:picname",func(c *gin.Context){
			picname:=c.Param("picname")
			picpath:=filepath.Join("public","2025pic",picname)
			data,err:=os.ReadFile(picpath)
			if err != nil{
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Data(http.StatusOK,"image/jpeg", data)
		})
	}

	apiGroup:=server.Group("/api")
	{
		apiGroup.GET("/playlist",func(c *gin.Context) {
			filePath := filepath.Join("public", "2024song", "playlist.json")
			data, err := os.ReadFile(filePath)
			if err != nil {
				c.String(http.StatusNotFound, "Playlist not found")
				return
			}
			c.Data(http.StatusOK, "application/json", data)
		})		
		apiGroup.GET("/songInfo", func(c *gin.Context) {
			songId := c.Query("songId")
			songUrl := fmt.Sprintf("https://music.163.com/api/v3/song/detail?c=[{\"id\":%s}]", songId)
			body, err := fetch(songUrl)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			c.Data(http.StatusOK, "application/json", body)
		})
		apiGroup.GET("/lyricsInfo", func(c *gin.Context) {
			songId := c.Query("songId")
			lyricsUrl := fmt.Sprintf("https://music.163.com/api/song/media?id=%s", songId)
			body, err := fetch(lyricsUrl)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			c.Data(http.StatusOK, "application/json", body)
		})
		apiGroup.GET("/picInfo", func(c *gin.Context) {
			url := c.Query("picUrl")
			body, err := fetch(url)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			c.Data(http.StatusOK, "image/jpeg", body)
		})
		apiGroup.GET("/flac", func(c *gin.Context) {
			songId := c.Query("songId")
			jsonPath := filepath.Join("public", "2024song", "playlist.json")
			data, _ := os.ReadFile(jsonPath)

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

					if _, err := os.Stat(mp3Path); err == nil {
						c.File(mp3Path)
						return
					} else if _, err := os.Stat(flacPath); err == nil {
						c.File(flacPath)
						return
					}
				}
			}

			c.String(http.StatusNotFound, "Song not found")
		})


	}

	server.Run(":4000")
}

