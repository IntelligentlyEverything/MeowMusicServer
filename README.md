- [English](#MeowMusicServer)
- [简体中文](#喵波音律)

# MeowMusicServer
[![codecov](https://codecov.io/gh/MoeCinnamo/MeowMusicServer/graph/badge.svg?token=20ZLNOK34R)](https://codecov.io/gh/MoeCinnamo/MeowMusicServer)

Your aggregated music API&private music player.

## Features
- Music API: Search, download, and stream music from various sources.(Under construction)

## Source code execution
``` sh
git clone https://github.com/xiaozhi-music/MeowMusicServer.git
cd MeowMusicServer
go mod tidy
cp .env.example .env # modify .env file according to your environment
mkdir music-uploads # create music upload directory
go run .
```

# 喵波音律
[![codecov](https://codecov.io/gh/MoeCinnamo/MeowMusicServer/graph/badge.svg?token=20ZLNOK34R)](https://codecov.io/gh/MoeCinnamo/MeowMusicServer)

你自己的音乐API、私人的音乐服务器

## 功能
- 音乐API：搜索、下载、直链播放。(开发中)

## 通过代码运行
``` sh
git clone https://github.com/xiaozhi-music/MeowMusicServer.git
cd MeowMusicServer
go mod tidy
cp .env.example .env # 根据你的环境修改 .env 文件
mkdir music-uploads # 创建音乐上载目录
go run .
```