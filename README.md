# iceapple
网易云音乐私人FM命令行版本

## 播放器控制
功能 | 快捷键
:--- | :---
播放/暂停 | Space
音量+ | Ctrl+Up
音量- | Ctrl+Down
快进 | S
快退 | s
下一首 | Ctrl+Right
喜欢/取消喜欢 | Ctrl+L
## 环境要求
+ glide
+ libgstreamer1.0
## 安装
```
go get github.com/yiitz/iceapple
cd $GOPATH/src/github.com/yiitz/iceapple
glide install
go install
```