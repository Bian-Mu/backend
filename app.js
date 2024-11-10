const express = require('express');
const cors = require('cors');
const fs = require('fs');
const path = require('path');
const app = express();
const port = 4000;

// 启用CORS
app.use(cors());

// 设置静态文件目录
app.use(express.static('public'));

// API端点：获取public/2024md目录下的指定Markdown文件
app.get('/public/2024md/:filename', (req, res) => {
    const filename = req.params.filename;
    const filePath = path.join(__dirname, 'public', '2024md', `${filename}`);

    fs.readFile(filePath, 'utf8', (err, data) => {
        if (err) {
            res.status(201).send('File not found');
            return;
        }
        res.send(data);
    });
});

// API端点：获取public/2024pic目录下的指定图片文件
app.get('/public/2024pic/:filename', (req, res) => {
    const filename = req.params.filename;
    const filePath = path.join(__dirname, 'public', '2024pic', `${filename}`);

    fs.readFile(filePath, (err, data) => {
        if (err) {
            res.status(201).send('File not found');
            return;
        }
        res.setHeader('Content-Type', `image/${path.extname(filename).slice(1)}`);
        res.send(data);
    });
});

//API端点：获取歌曲信息
app.get('/api/songInfo', async (req, res) => {
    const songId = req.query.songId;

    const songUrl = `https://music.163.com/api/v3/song/detail?c=[{"id":${songId}}]`;
    try {
        const responseInfo = await fetch(songUrl, {
            headers: {
                'Content-Type': 'application/json',
            }
        });

        const infoData = await responseInfo.json();
        res.json(infoData);

    } catch (error) {
        console.error('Error fetching data:', error);
        res.status(500).json({ error: 'Internal Server Error' });
    }
});

//API端点：获取歌词
app.get('/api/lyricsInfo', async (req, res) => {
    const songId = req.query.songId;

    const lyricsUrl = `https://music.163.com/api/song/media?id=${songId}`
    try {
        const responseLyrics = await fetch(lyricsUrl, {
            headers: {
                'Content-Type': 'application/json',
            }
        });
        const lyricsData = await responseLyrics.json();
        res.json(lyricsData)

    } catch (error) {
        console.error('Error fetching data:', error);
        res.status(500).json({ error: 'Internal Server Error' });
    }
});

//API端点：获取专辑封面
app.get("/api/picInfo", async (req, res) => {
    const url = req.query.picUrl;
    try {
        const responsePic = await fetch(url, {
            headers: {
                "Content-Type": "application/json"
            }
        })

        const imageBuffer = await responsePic.arrayBuffer();
        const buffer = Buffer.from(imageBuffer);
        res.set('Content-Type', 'image/jpg')
        res.send(buffer);
    } catch (error) {
        console.error('Error fetching data:', error);
        res.status(500).json({ error: 'Internal Server Error' });
    }
})

//API端点：获取歌单
app.get("/api/playlist", async (req, res) => {
    const filePath = path.join(__dirname, 'public', '2024song', 'playlist.json');
    fs.readFile(filePath, (err, data) => {
        if (err) {
            res.status(201).send('File not found');
            return;
        }
        const json = JSON.parse(data.toString())
        res.json(json);
    });
})

//API端点：获取flac
app.get("/api/flac", async (req, res) => {
    const songId = req.query.songId

    const jsonPath = path.join(__dirname, 'public', '2024song', 'playlist.json')

    const jsonData = JSON.parse(fs.readFileSync(jsonPath, 'utf8'));
    const songEntry = jsonData.find(entry => entry.id === songId.toString());

    if (!songEntry) {
        res.status(404).send('Song not found');
        return;
    }
    const filePath = path.join(__dirname, 'public', '2024song', 'music', `${songEntry.name}.flac`)

    res.sendFile(filePath, err => {
        if (err) {
            const filePath2 = path.join(__dirname, 'public', '2024song', 'music', `${songEntry.name}.mp3`)
            res.sendFile(
                filePath2, err2 => {
                    if (err2) {
                        res.status(201).send('File not found');
                        return;
                    }
                }
            )

        }
    })
})


app.listen(port, () => {
    console.log(`Server running at http://localhost:${port}`);
});
