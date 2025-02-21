package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

var PASSWORD string="lj050424"

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

func fileNameTest(fileName string) (int ,int,bool){
	regex := `^(\d+)-(\d+)\.md$`
	re := regexp.MustCompile(regex)

	matches := re.FindStringSubmatch(fileName)
	if len(matches) != 3 {
		return 0, 0, false
	}

	num1 := matches[1]
	num2 := matches[2]

	var month, day int
	_, err1 := fmt.Sscanf(num1, "%d", &month)
	_, err2 := fmt.Sscanf(num2, "%d", &day)

	if err1 != nil || err2 != nil {
		return 0, 0, false
	}

	return month, day, true
}


func saveFile(file io.Reader, fileName string) error {
	destPath := filepath.Join("./public","2025md",fileName)
	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	return err
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

		publicGroup.GET("/2024pic/:picname",func(c *gin.Context){
			picname:=c.Param("picname")
			picpath:=filepath.Join("public","2024pic",picname)
			data,err:=os.ReadFile(picpath)
			if err != nil{
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Data(http.StatusOK,"image/jpeg", data)
		})

		publicGroup.POST("uploadMD",func(c *gin.Context) {
			password := c.DefaultPostForm("password", "")
			if password!=PASSWORD{
				c.JSON(400, gin.H{
					"message": "Invalid file format. Expected format like 2-3.md",
				})
				return
			}


			file,header,_:=c.Request.FormFile(("file"))
			name:=header.Filename
			defer file.Close()

			month,day,result:=fileNameTest(name)
			if !result{
				c.JSON(400, gin.H{
					"message": "expect filename likes 2-3.md",
				})
				return
			}else if month<13 && month>0{
				currentTime:=time.Now()
				firstDay:=time.Date(currentTime.Year(),time.Month(month),1,0,0,0,0,currentTime.Location())
				nextMonth := firstDay.AddDate(0, 1, 0)	
				daysInMonth := nextMonth.Sub(firstDay).Hours() / 24

				if day<=int(daysInMonth) && day>0{
					err:=saveFile(file,name)
					if err==nil{
						c.JSON(200, gin.H{
							"message": "upload successfully",
						})
					}
				}else{
					c.JSON(400, gin.H{
						"message": "day is incorrect",
					})
				return
				}
			}else{
				c.JSON(400, gin.H{
					"message": "month is incorrect",
				})
				return
			}
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

