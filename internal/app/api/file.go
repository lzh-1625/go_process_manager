package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/logic"

	"github.com/gin-gonic/gin"
)

type file struct{}

var FileApi = new(file)

func (f *file) FilePathHandler(ctx *gin.Context) {
	path := getQueryString(ctx, "path")
	rOk(ctx, "Operation successful!", logic.FileLogic.GetFileAndDirByPath(path))
}

func (f *file) FileWriteHandler(ctx *gin.Context) {
	path := ctx.PostForm("filePath")
	fi, err := ctx.FormFile("data")
	errCheck(ctx, err != nil, "Read file data failed!")
	fiReader, _ := fi.Open()
	err = logic.FileLogic.UpdateFileData(path, fiReader, fi.Size)
	errCheck(ctx, err != nil, "Update file data operation failed!")
	rOk(ctx, "Operation successful!", nil)
}

func (f *file) FileReadHandler(ctx *gin.Context) {
	path := getQueryString(ctx, "filePath")
	bytes, err := logic.FileLogic.ReadFileFromPath(path)
	errCheck(ctx, err != nil, "Operation failed!")
	rOk(ctx, "Operation successful!", string(bytes))
}
