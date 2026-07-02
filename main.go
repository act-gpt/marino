package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/act-gpt/marino/api"
	"github.com/act-gpt/marino/common"
	r "github.com/act-gpt/marino/common/redis"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/middleware"
	"github.com/act-gpt/marino/model"
	rt "github.com/act-gpt/marino/router"
	"github.com/act-gpt/marino/web"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	stats "github.com/semihalev/gin-stats"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
)

var server *http.Server

func handleDebug() {
	c := make(chan os.Signal, 1) // channel to wait os.Signal
	signal.Notify(c, syscall.SIGHUP)
	go func() { // go routine to do not block other requests on the application
		<-c
		system.Config.ShowPrompt = !system.Config.ShowPrompt
		fmt.Println("Handel signal, ShowPrompt:", system.Config.ShowPrompt)
	}()
}

func handleShutdown() {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func main() {

	version := fmt.Sprintf("ACT GPT/%s (%s %s)", config.Version, runtime.GOARCH, runtime.GOOS)
	fmt.Println("\033[32;1;4m" + version + "\033[0m")
	godotenv.Load("./etc/.env")
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if gin.Mode() == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	source := os.Getenv("DB_DATA_SOURCE")
	if err := model.InitDB(source); err != nil {
		fmt.Println("\033[31;1;4mfailed to initialize database: " + err.Error() + "\033[0m")
		panic(err)
	}

	conf, _ := model.LoadSystemConfig()
	if save := system.InitNeedSave(conf); save {
		fmt.Print(system.Config)
		/*
			model.InsertOrSaveSystemConfig(system.Config)
			//TODO: Initialled 判断
			conf = system.Config
		*/
	}

	mode := "console"
	if gin.Mode() != "debug" {
		gin.SetMode(gin.ReleaseMode)
		mode = "file"
	}

	logx.DisableStat()
	logx.SetUp(logx.LogConf{
		ServiceName: system.Config.SystemName,
		Mode:        mode,
		KeepDays:    30,
	})

	api.NewApiClient()

	proc.AddShutdownListener(func() {
		err := model.CloseDB()
		if err != nil {
			fmt.Println("\033[31;1;4mfailed to close database: " + err.Error() + "\033[0m")
		}
		r.Close()
	})

	// Initialize HTTP server
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(stats.RequestStats())
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	if conf.Initialled.Redis {
		r.Init()
		opt := r.ParseRedisOption()
		store, _ := redis.NewStore(opt.MinIdleConns, opt.Network, opt.Addr, opt.Password, []byte(conf.SessionSecret))
		router.Use(sessions.Sessions("SESSION", store))
	} else {
		store := cookie.NewStore([]byte(conf.SessionSecret))
		router.Use(sessions.Sessions("SESSION", store))
	}

	rt.SetRouter(router, web.BuildFS)

	var host = conf.Host
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(conf.Port)
	}

	server = &http.Server{
		Addr:    host + ":" + port,
		Handler: router,
	}

	go func() {
		time.Sleep(1 * time.Second)
		url := os.Getenv("PORTAL")
		if url != "" {
			if !strings.HasPrefix(url, "http") {
				url += "http://" + host + ":" + port + url
			}
			common.Open(url)
		}
	}()
	go func() {
		fmt.Println("\033[32;1;4mServer start listening at " + host + ":" + port + "\033[0m")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		} else {
			fmt.Println("\033[31;1;4m Started \033[0m")
		}
	}()
	//handleDebug()
	handleShutdown()

	/*
		err := server.Run(host + ":" + port)
		if err != nil {
			fmt.Println("\033[31;1;4mfailed to start HTTP server: " + err.Error() + "\033[0m")
		}
	*/
}
