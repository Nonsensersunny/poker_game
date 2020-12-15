package main

import (
	"flag"
	"github.com/Nonsensersunny/poker_game/pkg/client"
	"github.com/Nonsensersunny/poker_game/pkg/controller"
	"github.com/Nonsensersunny/poker_game/pkg/server"
	"github.com/Nonsensersunny/poker_game/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"net/http"
	"os"
)

func initLogger() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&prefixed.TextFormatter{
		DisableColors: false,
		TimestampFormat : "2006-01-02 15:04:05",
		FullTimestamp:true,
		ForceFormatting: true,
	})
}

func init() {
	initLogger()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func initHTTP() {
	r := gin.Default()
	//corsConfig := cors.DefaultConfig()
	//corsConfig.AllowCredentials = true
	////corsConfig.AllowOrigins = []string{"http://10.108.15.55:8080"}
	//corsConfig.AllowAllOrigins = true
	//r.Use(cors.New(corsConfig))
	r.Use(Cors())

	defer r.Run(":8888")

	var (
		gameGroup = r.Group("/game")
		userGroup = r.Group("/user")
	)

	// game
	{
		gameGroup.GET("/", controller.GetAvailableGames)
		gameGroup.POST("/", controller.InitGame)
		gameGroup.PUT("/", controller.StartGame)
		gameGroup.DELETE("/", controller.DestroyGame)
		gameGroup.POST("/remain", controller.TakeRemain)
		gameGroup.GET("/remain", controller.UncoverRemain)
	}

	// user
	{
		userGroup.GET("/", controller.CheckName)
		userGroup.POST("/game", controller.JoinGame)
		userGroup.GET("/game", controller.GetPlayer)
		userGroup.PUT("/game", controller.Play)
	}
}

func main()  {
	//go initHTTP()
	var (
		service = flag.String("service", util.ServiceClient, "service")
		address = flag.String("address", "localhost:8080", "address")
	)
	flag.Parse()

	if service != nil {
		if *service == util.ServiceServer {
			server.InitServer(*address)
		} else {
			client.InitClient(*address)
		}
	}


	//g := game.NewGame(game.ModeChinesePoker)
	//g.Shuffle()
	////for _, v := range game.Cards {
	////	fmt.Printf("%v %v\n", v.Name, v.Color.ToUnicode())
	////}
	//deals := g.Deal()
	//deals[0].Sort()
	//for _, v := range deals[0] {
	//	fmt.Printf("%v %v\n", v.Name, v.Color.ToUnicode())
	//}
	//client.InitUI(deals[0], model.Players{
	//	{
	//		Name: "player1",
	//		Play: deals[1],
	//	},
	//	{
	//		Name: "player2",
	//		Play: deals[2],
	//	},
	//})
}
