package cmd

import (
	"context"
	"fmt"
	convert "godoai/convert"
	knowledge "godoai/knowledge"
	model "godoai/model"
	progress "godoai/progress"
	sd "godoai/sd"
	"godoai/serv"
	sys "godoai/sys"
	voice "godoai/voice"
	"log"
	"net/http"
	"time"
)

var server = serv.GetServer(":56710")

func Start() error {
	err := InitSystem()
	if err != nil {
		return fmt.Errorf("init system error: %v", err)
	}
	// 注册路由
	server.RegisterRoutes(func(router *serv.RestrictedRouter) {
		router.Handle("/system/setting", http.MethodPost, sys.ConfigHandle)
		router.Handle("/system/updateInfo", http.MethodGet, sys.GetUpdateUrlHandler)
		router.Handle("/system/update", http.MethodGet, sys.UpdateAppHandler)

		router.Handle("/knowledge/list", http.MethodGet, knowledge.ListHandle)
		router.Handle("/knowledge/create", http.MethodPost, knowledge.CreateHandle)
		router.Handle("/knowledge/delete", http.MethodPost, knowledge.DeleteHandle)
		router.Handle("/knowledge/deleteFile", http.MethodPost, knowledge.DeleteFileHandle)
		router.Handle("/knowledge/ask", http.MethodPost, knowledge.AskHandle)
		router.Handle("/knowledge/add", http.MethodPost, knowledge.AddHandle)
		router.Handle("/knowledge/upload", http.MethodPost, convert.MultiUploadHandler)
		router.Handle("/knowledge/url", http.MethodPost, convert.HandleURLPost)
		router.Handle("/knowledge/filedetail", http.MethodGet, convert.ShowDetailHandler)
		router.Handle("/knowledge/showimage", http.MethodGet, convert.ServeImage)
		router.Handle("/knowledge/convertfile", http.MethodPost, convert.JsonParamHandler)
		router.Handle("/model/download", http.MethodPost, model.Download)
		router.Handle("/model/outserver", http.MethodGet, model.DownServerHandler)
		router.Handle("/model/delete", http.MethodPost, model.DeleteFileHandle)
		router.Handle("/model/tags", http.MethodGet, model.Tagshandler)
		router.Handle("/model/show", http.MethodGet, model.ShowHandler)
		router.Handle("/model/chat", http.MethodPost, model.ChatHandler)
		router.Handle("/model/embeddings", http.MethodPost, model.EmbeddingHandler)
		router.Handle("/model/uploadimage", http.MethodPost, sd.UploadHandler)
		router.Handle("/model/image", http.MethodPost, sd.CreateImage)
		router.Handle("/model/deleteimage", http.MethodPost, sd.DeleteImageHandler)
		router.Handle("/model/viewimage", http.MethodGet, sd.ServeImage)
		router.Handle("/model/voice", http.MethodPost, voice.UploadHandler)
		router.Handle("/model/tts", http.MethodPost, voice.TtsHandler)
		router.Handle("/model/audio", http.MethodGet, voice.ServeAudio)

	})
	server.Start()
	return nil
}
func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := progress.StopAllCmd()
	if err != nil {
		log.Fatalf("Servers forced to shutdown error: %v", err)
	}
	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped.")
}
