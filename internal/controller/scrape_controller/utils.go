package scrape_controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas"
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path"
	"path/filepath"
)

// 保存图片到指定路径
// imgPath: 图片保存路径
// img: 要保存的图片
func SaveImage(imgFile *schemas.FileInfo, img image.Image) error {
	var (
		buff bytes.Buffer
		err  error
	)
	switch filepath.Ext(imgFile.Name()) {
	case ".jpg", ".jpeg": // 保存为 JPEG
		err = png.Encode(&buff, img)
	default: // 保存为 PNG
		err = jpeg.Encode(&buff, img, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return err
	}
	return storage_controller.CreateFile(imgFile, &buff)
}

// 下载 TMDB 图片并保存到指定路径
// 自动根据 TMDB 图片的扩展名决定保存格式
// p: TMDB 中图片地址
// target: 目标路径，不带后缀名
func DownloadTMDBImageAndSave(p string, target string, storageType schemas.StorageType) error {
	target += path.Ext(p)
	dstFile, err := storage_controller.GetFile(target, storageType)
	if err != nil {
		return err
	}
	exists, err := storage_controller.Exists(dstFile)
	if err != nil {
		return err
	}
	if exists {
		return nil // 如果文件已存在，则跳过下载
	}

	img, err := tmdb_controller.DownloadImage(p)
	if err != nil {
		return err
	}
	return SaveImage(dstFile, img)
}

// 下载 Fanart 图片并保存到指定路径
// 自动根据 Fanart 图片的扩展名决定保存格式
// url: Fanart 中图片地址
// target: 目标路径，不带后缀名
func DownloadFanartImageAndSave(url string, target string, storageType schemas.StorageType) error {
	target += path.Ext(url)
	dstFile, err := storage_controller.GetFile(target, storageType)
	if err != nil {
		return err
	}
	exists, err := storage_controller.Exists(dstFile)
	if err != nil {
		return err
	}
	if exists {
		return nil // 如果文件已存在，则不下载
	}

	img, err := fanart_controller.DownloadImage(url)
	if err != nil {
		return err
	}
	return SaveImage(dstFile, img)
}

func bytes2Reader(p []byte) (io.Reader, error) {
	var buffer bytes.Buffer
	_, err := buffer.Write(p)
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}
