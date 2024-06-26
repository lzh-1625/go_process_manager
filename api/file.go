package api

import (
	FileService "msm/service/file"

	"github.com/gin-gonic/gin"
)

type file struct{}

var FileApi = new(file)

func (f *file) FilePathHandler(ctx *gin.Context) {
	data, err := FileService.FileService.GetFileAndDirByPath(ctx.Query("path"))
	errCheck(ctx, err != nil, "文件路径查询失败")
	rOk(ctx, "文件路径查询成功", data)
}

func (f *file) FileWriteHandler(ctx *gin.Context) {
	path := ctx.PostForm("filePath")
	fi, err := ctx.FormFile("data")
	errCheck(ctx, err != nil, "文件读取失败")
	fiReader, _ := fi.Open()
	err = FileService.FileService.UpdateFileData(path, fiReader, fi.Size)
	errCheck(ctx, err != nil, "文件数据更新失败")
	rOk(ctx, "文件更新成功", nil)
}

func (f *file) FileReadHandler(ctx *gin.Context) {
	path := ctx.Query("filePath")
	bytes, err := FileService.FileService.ReadFileFromPath(path)
	errCheck(ctx, err != nil, "文件数据读取失败")
	rOk(ctx, "文件数据读取成功", string(bytes))
}
