# 最伟大的作品

最快的 JableTV 下载器

## 使用

```shell
Usage of app:
    -c int
        concurrent download num range 0-30 (default 5)
    -u string
        input your url
```

## 特性

- rod 无头浏览器
- 不使用 ffmpeg
- 使用 gomedia 处理视频

## benchmark

网络: 300M 带宽

设备: mac m1 air

| 本仓库用时 | hcjohn463 用时 |
|-------|--------------|
| 1m34s | 11m54s       |

## 鸣谢

- aes 加解密参考 https://github.com/hcjohn463/JableTVDownload
- 合并 mp4 参考 https://github.com/orestonce/m3u8d/blob/main/merge.go