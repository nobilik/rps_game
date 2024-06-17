package db

import (
	"fmt"
	"rps/models"

	"os"

	"gorm.io/driver/mysql"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	// DB *sql.DB
	DB          *gorm.DB
	RedisClient *redis.Client
)

// DBConnect here we access db - to init
func SetDB() {

	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	// Формирование строки подключения
	dsn := fmt.Sprintf("%s:%s@tcp(dbContainer)/%s?parseTime=true", dbUser, dbPassword, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to MySQL database!")

	migrate()
}

func migrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Bet{},
		&models.Transaction{},
	)

}

func SetRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis" // имя сервиса из docker-compose.yml
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379" // стандартный порт Redis
	}

	addr := redisHost + ":" + redisPort

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // укажите пароль, если требуется
		DB:       0,  // используем стандартную базу данных
	})

	// // Проверка соединения с сервером Redis
	// _, err := RedisClient.Ping().Result()
	// if err != nil {
	// 	log.Fatalf("Не удалось подключиться к Redis: %v", err)
	// }

	// log.Println("Подключение к Redis установлено")
}
