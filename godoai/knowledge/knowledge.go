package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"godoai/knowledge/dbtype"
	"godoai/knowledge/embeddings"
	"godoai/libs"
)

// newFactory 创建数据库工厂，并处理可能的错误。
func newFactory(config dbtype.DbConfig) (*dbtype.DbFactory, error) {
	factory, err := NewDbFactory(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create database: %w", err)
	}
	return factory, nil
}

// decodeRequest 解析请求体并返回请求参数，如果发生错误，则调用handleError。
func decodeRequest(w http.ResponseWriter, r *http.Request, request interface{}) error {
	//err := json.NewDecoder(r.Body).Decode(request)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		libs.Error(w, err.Error())
		return err
	}
	return nil
}

func CreateHandle(w http.ResponseWriter, r *http.Request) {
	var request dbtype.CreateParams
	if err := decodeRequest(w, r, &request); err != nil {
		return
	}
	db, err := newFactory(request.Config)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	res, err := db.DB.Create(request.Name, request.Model)

	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	libs.Success(w, res, "Create Success")
}

func ListHandle(w http.ResponseWriter, r *http.Request) {
	var request dbtype.ConfigParams
	if err := decodeRequest(w, r, &request); err != nil {
		return
	}
	db, err := newFactory(request.Config)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	res, err := db.DB.List()

	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	libs.Success(w, res, "success")
}

func DeleteHandle(w http.ResponseWriter, r *http.Request) {
	var request dbtype.DeleteParams
	if err := decodeRequest(w, r, &request); err != nil {
		return
	}
	db, err := newFactory(request.Config)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	err = db.DB.Delete(request.Name)

	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	libs.Success(w, "", "success")
}

func AddHandle(w http.ResponseWriter, r *http.Request) {
	var request dbtype.AddParams
	if err := decodeRequest(w, r, &request); err != nil {
		return
	}
	db, err := newFactory(request.Config)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	err = db.DB.Add(request.Name, request.Model, request.Files)

	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	libs.Success(w, "", "success")
}

func AskHandle(w http.ResponseWriter, r *http.Request) {
	var request dbtype.AskParams
	if err := decodeRequest(w, r, &request); err != nil {
		return
	}
	//log.Printf("Ask: %+v", request)
	db, err := newFactory(request.Config)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	res, err := db.DB.Ask(request.Name, request.Model, request.Message)

	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	res = embeddings.SelectDocuments(res)
	libs.Success(w, res, "success")
}

func DeleteFileHandle(w http.ResponseWriter, r *http.Request) {
	var request dbtype.DeleteFilearams
	if err := decodeRequest(w, r, &request); err != nil {
		return
	}
	db, err := newFactory(request.Config)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	err = db.DB.DeleteFile(request.Name, request.Model, request.File)

	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	// 删除文件
	err = os.Remove(request.File)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	err = os.Remove(request.File + "_result.json")
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	if libs.PathExists(request.File + "_images") {
		err = os.RemoveAll(request.File + "_images")
		if err != nil {
			libs.Error(w, err.Error())
			return
		}
	}
	libs.Success(w, "", "success")
}
