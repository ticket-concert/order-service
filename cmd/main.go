package main

import (
	"fmt"
	logGo "log"
	"order-service/configs"
	eventRepoQuery "order-service/internal/modules/event/repositories/queries"
	orderHandler "order-service/internal/modules/order/handlers"
	orderRepoCommand "order-service/internal/modules/order/repositories/commands"
	orderRepoQuery "order-service/internal/modules/order/repositories/queries"
	orderUsecase "order-service/internal/modules/order/usecases"
	roomHandler "order-service/internal/modules/room/handlers"
	roomRepoCommand "order-service/internal/modules/room/repositories/commands"
	roomRepoQuery "order-service/internal/modules/room/repositories/queries"
	roomUsecase "order-service/internal/modules/room/usecases"
	ticketRepoCommand "order-service/internal/modules/ticket/repositories/commands"
	ticketRepoQuery "order-service/internal/modules/ticket/repositories/queries"
	userRepoQuery "order-service/internal/modules/user/repositories/queries"
	"order-service/internal/pkg/apm"
	"order-service/internal/pkg/databases/mongodb"
	graceful "order-service/internal/pkg/gs"
	"order-service/internal/pkg/helpers"
	kafkaConfluent "order-service/internal/pkg/kafka/confluent"
	"order-service/internal/pkg/log"
	"order-service/internal/pkg/redis"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.elastic.co/apm/module/apmfiber"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// @BasePath	/
func main() {
	// Init Config
	configs.InitConfig()

	// Init Elastic APM Config or DD Apm
	if configs.GetConfig().Datadog.DatadogEnabled == "true" {
		tracer.Start(
			tracer.WithAgentAddr(fmt.Sprintf("%s:%s", configs.GetConfig().Datadog.DatadogHost, configs.GetConfig().Datadog.DatadogPort)),
			tracer.WithEnv(configs.GetConfig().Datadog.DatadogEnv),
			tracer.WithService(configs.GetConfig().Datadog.DatadogService),
		)
		defer tracer.Stop()
	} else {
		// temp handling until we move all to Datadog
		apm.InitConnection(configs.GetConfig().APMElastic.APMUrl, configs.GetConfig().APMElastic.APMSecretToken)
	}
	// Init MongoDB Connection
	mongo := mongodb.MongoImpl{}
	mongo.SetCollections(&mongo)
	mongo.InitConnection(configs.GetConfig().MongoDB.MongoMasterDBUrl, configs.GetConfig().MongoDB.MongoSlaveDBUrl)
	// Init Logger
	logZap := log.SetupLogger(configs.GetConfig().ServiceName)
	log.Init(logZap)

	// Init BlacklistedEmail
	helpers.InitReadBlackListEmail()

	// Init Kafka Config
	kafkaConfluent.InitKafkaConfig(configs.GetConfig().Kafka.KafkaUrl, configs.GetConfig().Kafka.KafkaUsername, configs.GetConfig().Kafka.KafkaPassword)

	// Init instance fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})
	app.Use(apmfiber.Middleware(apmfiber.WithTracer(apm.GetTracer())))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		// Format:       `${time} {"router_activity" : [${status},"${latency}","${method}","${path}"], "query_param":${queryParams}, "body_param":${body}}` + "\n",
		TimeInterval: time.Millisecond,
		TimeFormat:   "02-01-2006 15:04:05",
		TimeZone:     "Indonesia/Jakarta",
	}))
	shutdownDelay, _ := strconv.Atoi(configs.GetConfig().ShutDownDelay)
	// graceful shutdown setup
	gs := &graceful.GracefulShutdown{
		Timeout:        5 * time.Second,
		GracefulPeriod: time.Duration(shutdownDelay) * time.Second,
	}
	app.Get("/healthz", gs.LivenessCheck)
	app.Get("/readyz", gs.ReadinessCheck)
	gs.Enable(app)

	setHttp(app, gs)

	//=== listen port ===//
	if err := app.Listen(fmt.Sprintf(":%s", configs.GetConfig().ServicePort)); err != nil {
		logGo.Fatal(err)
	}
}

func setHttp(app *fiber.App, gs *graceful.GracefulShutdown) {
	// Init Redis
	redisClient := redis.InitConnection(configs.GetConfig().Redis.RedisDB, configs.GetConfig().Redis.RedisHost, configs.GetConfig().Redis.RedisPort,
		configs.GetConfig().Redis.RedisPassword, configs.GetConfig().Redis.RedisAppConfig)
	// Init Jwt
	helperImpl := &helpers.JwtImpl{}
	helperImpl.InitConfig(configs.GetConfig().Jwt.JwtPrivateKey, configs.GetConfig().Jwt.JwtPublicKey,
		configs.GetConfig().Jwt.JwtRefreshPrivateKey, configs.GetConfig().Jwt.JwtRefreshPublicKey)

	logger := log.GetLogger()
	mongoMasterClient := mongodb.NewMongoDBLogger(mongodb.GetMasterConn(), mongodb.GetMasterDBName(), logger)
	mongoSlaveClient := mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetMasterDBName(), logger)
	kafkaProducer, err := kafkaConfluent.NewProducer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), logger)
	if err != nil {
		panic(err)
	}
	gs.Register(
		mongoMasterClient,
		mongoSlaveClient,
		graceful.FnWithError(redisClient.Close),
		kafkaProducer,
	)

	ticketQueryMongodbRepo := ticketRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)
	ticketCommandMongodbRepo := ticketRepoCommand.NewCommandMongodbRepository(mongoMasterClient, logger)

	eventQueryMongodbRepo := eventRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)
	userQueryMongodbRepo := userRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)

	roomCommandMongodbRepo := roomRepoCommand.NewCommandMongodbRepository(mongoMasterClient, logger)
	roomQueryMongodbRepo := roomRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)
	roomUsecase := roomUsecase.NewCommandUsecase(roomQueryMongodbRepo, roomCommandMongodbRepo, ticketQueryMongodbRepo,
		eventQueryMongodbRepo, logger, redisClient)

	orderCommandMongodbRepo := orderRepoCommand.NewCommandMongodbRepository(mongoMasterClient, logger)
	orderQueryMongodbRepo := orderRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)
	orderUsecaseCommand := orderUsecase.NewCommandUsecase(orderCommandMongodbRepo, orderQueryMongodbRepo, roomQueryMongodbRepo,
		ticketQueryMongodbRepo, ticketCommandMongodbRepo, eventQueryMongodbRepo, userQueryMongodbRepo, logger, redisClient)
	orderUsecaseQuery := orderUsecase.NewQueryUsecase(orderQueryMongodbRepo, logger)

	// set module
	roomHandler.InitRoomHttpHandler(app, roomUsecase, logger, redisClient)
	orderHandler.InitOrderHttpHandler(app, orderUsecaseCommand, orderUsecaseQuery, logger, redisClient)

}
