# 无人直播111
![GitHub License](https://img.shields.io/github/license/无人直播/无人直播)
![GitHub Release](https://img.shields.io/github/v/release/无人直播/无人直播)
![GitHub stars](https://img.shields.io/github/stars/无人直播/无人直播?style=social)
![GitHub forks](https://img.shields.io/github/forks/无人直播/无人直播?style=social)
![GitHub issues](https://img.shields.io/github/issues/无人直播/无人直播)
![GitHub last commit](https://img.shields.io/github/last-commit/无人直播/无人直播)
![GitHub all releases](https://img.shields.io/github/downloads/无人直播/无人直播/total)
![Docker Pulls](https://img.shields.io/docker/pulls/无人直播/无人直播)


[简体中文](README.zh_cn.md) | [Official Website](https://www.无人直播.com/) | [Documentation](https://docs.无人直播.com/)

无人直播是一款简单好用的AI助手，支持知识库、模型API、分享、联网搜索、智能体，它还在飞快成长中。

无人直播 is an easy-to-use AI assistant that supports knowledge bases, model APIs, sharing, web search, and intelligent agents. It's rapidly growing and improving.

## 🚀 One-sentence Introduction  

A user-friendly AI assistant software that supports local AI models, APIs, and knowledge base setup.

## ✅ Core Features  

- One-click deployment of local AI models and mainstream model APIs
![Local model](.github/assets/img/1_en.png)
- Local knowledge base
![Knowledge base](.github/assets/img/3_en.png)
- Intelligent agent creation
![Intelligent agent](.github/assets/img/4_en.png)
  
- Can be shared online for others to use
![Sharing](.github/assets/img/5_en.png)

- Supports web search
![Web search](.github/assets/img/6_en.png)

- Supports server-side deployment 

- MCP Client
![MCP Client](.github/assets/img/7_en.png)

- Simultaneous conversations with multiple models in a single session (coming soon)  

## ✨ Product Highlights  
- Simple and easy to use, beginner-friendly for AI newcomers  

## 📥 Quick Installation

### Client Version（MacOS, Windows） 

- [Download from official website](https://www.无人直播.com/)   
- [Download from CNB](https://cnb.cool/无人直播/无人直播/-/releases/)  
- [Download from Github](https://github.com/无人直播/无人直播/releases)  

### Server Version

#### Docker Run
```bash 
docker run -d \
  --name node \
  -v $(pwd)/data:/无人直播/data \
  -v $(pwd)/uploads:/无人直播/uploads \
  -v $(pwd)/logs:/无人直播/logs \
  -p 7071:7071 \
  -w /无人直播 \
  无人直播/无人直播
```

#### Docker Compose
```bash
mkdir -p 无人直播
cd 无人直播
wget https://cnb.cool/无人直播/无人直播/-/git/raw/server/docker-compose.yml
# Run
docker compose up -d
# or
docker-compose up -d
``` 
## Build
```bash
git clone https://github.com/无人直播/无人直播.git
cd 无人直播
# For macOS users, please remove the `@rollup/rollup-win32-x64-msvc` dependency in [package.json](http://_vscodecontentref_/0)
cd frontend
yarn
cd ..
yarn
yarn dev
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=无人直播/无人直播&type=Date)](https://www.star-history.com/#无人直播/无人直播&Date)