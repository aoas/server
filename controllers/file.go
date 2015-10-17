package controllers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aoas/server/utils"

	"github.com/aoas/server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type File struct {
	Base
}

func init() {
	addPermission("file", "file.list", "查询上传列表")
	addPermission("file", "file.upload", "上传文件")

}

func (f *File) Find(c *gin.Context) {
	if !isGranted(c, "file.list") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}

	page := queryInt(c, "page", 1)
	pagesize := queryInt(c, "pagesize", 30)
	userid := queryInt64(c, "user_id", 0)
	ext := c.Query("ext")
	name := c.Query("name")

	engine := models.Engine()
	session := func() *xorm.Session {
		where := engine.Where("")
		if userid > 0 {
			where = where.Where("user_id = ?", userid)
		}

		if ext != "" {
			exts := strings.Split(ext, ",")
			where = where.In("ext", exts)
		}

		if name != "" {
			where = where.Where("name like ?", "%"+name+"%")
		}

		return where
	}

	total, _ := session().Count(&models.File{})
	files := make([]models.File, 0)

	if err := session().Limit(pagesize, pagesize*(page-1)).Desc("created_at").Find(&files); err != nil {
		c.JSON(400, utils.NewError("find files failed - %s", err.Error()))
		return
	}

	result := models.NewQueryResult(page, pagesize, total, files)

	c.JSON(200, result)

}

func (f *File) Upload(c *gin.Context) {
	if !isGranted(c, "file.upload") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	cu := currentUser(c)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, utils.NewError("read file content failed - %s", err.Error()))
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%v%v%s", cu.Id, time.Now().UnixNano(), ext)
	filefolder := "/" + time.Now().Format("20060102")
	path := f.Config.File.UploadPath + filefolder

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(400, utils.NewError("read file content failed - %s", err.Error()))
		return
	}

	fileModel := models.File{}
	fileModel.Key = filename
	fileModel.Name = header.Filename
	fileModel.Path = filefolder
	fileModel.Size = len(buf)
	fileModel.UserId = cu.Id
	fileModel.Ext = ext

	if err := fileModel.CheckValid(); err != nil {
		c.JSON(400, utils.NewError("upload file failed - %s", err.Error()))
		return
	}

	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(path, os.ModePerm)
	}

	fmt.Println("file path:", path)

	x := models.Engine()
	s := x.NewSession()
	s.Begin()

	if _, err := s.Insert(&fileModel); err != nil {
		c.JSON(400, utils.NewError("write database failed - %s", err.Error()))
		return
	}

	if err := ioutil.WriteFile(path+"/"+filename, buf, os.ModePerm); err != nil {
		s.Rollback()
		c.JSON(400, utils.NewError("write file failed - %s", err.Error()))
		return
	}

	s.Commit()

	c.JSON(200, fileModel)

}
func (f *File) Delete(c *gin.Context) {

}
func (f *File) ResizeImage(c *gin.Context) {
}
