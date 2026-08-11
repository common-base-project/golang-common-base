package connection

import (
	"context"
	"fmt"
	"golang-common-base/pkg/config"
	"golang-common-base/pkg/logger"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database struct {
	Self *gorm.DB
}

var DB *Database
var RedisClient *redis.Client
var MinioClient *minio.Client
var CasbinEnforcer *casbin.Enforcer
var OIDCVerifier *oidc.IDTokenVerifier

func Initial() {
	DB.Init()
	InitRedis()
	InitMinio()
	InitAuthz()
}

func openMySQLDB(username, password, addr, name, dsn string) *gorm.DB {
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
			username,
			password,
			addr,
			name,
			true,
			"Local")
	}

	logger.Info("database mysql dsn configured")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         128,
		SkipInitializeWithVersion: false,
	}), buildGormConfig())
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败，连接地址: %s，error: %s", addr, err))
	}
	setupPool(db)
	return db
}

func openSQLiteDB(path string) *gorm.DB {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		panic(fmt.Sprintf("创建 sqlite 目录失败: %v", err))
	}

	db, err := gorm.Open(sqlite.Open(path), buildGormConfig())
	if err != nil {
		panic(fmt.Sprintf("sqlite 连接失败，路径: %s，error: %s", path, err))
	}
	setupPool(db)
	return db
}

func buildGormConfig() *gorm.Config {
	return &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	}
}

func setupPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库连接池失败: %v", err))
	}
	sqlDB.SetMaxOpenConns(viper.GetInt(`db.gorm.maxOpenConn`))
	sqlDB.SetMaxIdleConns(viper.GetInt(`db.gorm.maxIdleConn`))
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func openDB() *gorm.DB {
	if config.DBConfig.Driver == "sqlite" {
		return openSQLiteDB(config.DBConfig.SQLite)
	}
	return openMySQLDB(
		config.DBConfig.UserName,
		config.DBConfig.Password,
		config.DBConfig.URL,
		config.DBConfig.DBName,
		config.DBConfig.DSN,
	)
}

func (db *Database) Init() {
	DB = &Database{
		Self: openDB(),
	}
}

func InitRedis() {
	addr := strings.TrimSpace(viper.GetString("redis.addr"))
	if addr == "" {
		logger.Warn("redis 未配置，跳过初始化")
		return
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		logger.Errorf("redis 初始化失败: %v", err)
		RedisClient = nil
		return
	}
	logger.Info("redis 初始化完成")
}

func InitMinio() {
	endpoint := strings.TrimSpace(viper.GetString("minio.endpoint"))
	endpoints := viper.GetStringSlice("minio.endpoints")
	if endpoint == "" && len(endpoints) > 0 {
		endpoint = strings.TrimSpace(endpoints[0])
	}
	if endpoint == "" {
		logger.Warn("minio 未配置，跳过初始化")
		return
	}

	useSSL := viper.GetBool("minio.use_ssl")
	endpoint, useSSL = normalizeMinioEndpoint(endpoint, useSSL)
	accessKey := viper.GetString("minio.access_key")
	secretKey := viper.GetString("minio.secret_key")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		for _, ep := range endpoints {
			ep = strings.TrimSpace(ep)
			if ep == "" || ep == endpoint {
				continue
			}
			ep, useSSL = normalizeMinioEndpoint(ep, useSSL)
			client, err = minio.New(ep, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
				Secure: useSSL,
			})
			if err == nil {
				endpoint = ep
				break
			}
		}
		if err != nil {
			logger.Errorf("minio 初始化失败: %v", err)
			return
		}
	}

	MinioClient = client
	logger.Infof("minio 初始化完成, endpoint=%s", endpoint)
}

func normalizeMinioEndpoint(raw string, fallbackSSL bool) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw, fallbackSSL
	}

	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			return raw, fallbackSSL
		}
		return u.Host, u.Scheme == "https"
	}

	return raw, fallbackSSL
}

func InitAuthz() {
	if DB == nil || DB.Self == nil {
		logger.Warn("数据库未初始化，跳过 Casbin 初始化")
		return
	}

	modelPath := strings.TrimSpace(viper.GetString("authz.casbin.model_path"))
	if modelPath == "" {
		modelPath = "./conf/casbin_model.conf"
	}

	adapter, err := gormadapter.NewAdapterByDB(DB.Self)
	if err != nil {
		logger.Errorf("Casbin adapter 初始化失败: %v", err)
		return
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		logger.Errorf("Casbin enforcer 初始化失败: %v", err)
		return
	}

	if err = enforcer.LoadPolicy(); err != nil {
		logger.Errorf("Casbin 策略加载失败: %v", err)
		return
	}

	if viper.GetBool("authz.casbin.bootstrap.enabled") {
		subject := strings.TrimSpace(viper.GetString("authz.casbin.bootstrap.subject"))
		if subject != "" {
			has, _ := enforcer.HasPolicy(subject, "/*", "(get|post|put|patch|delete|head|options)")
			if !has {
				_, _ = enforcer.AddPolicy(subject, "/*", "(get|post|put|patch|delete|head|options)")
				_ = enforcer.SavePolicy()
			}
		}
	}

	CasbinEnforcer = enforcer
	logger.Info("Casbin 初始化完成")

	if !viper.GetBool("auth.casdoor.enabled") {
		return
	}

	issuer := strings.TrimSpace(viper.GetString("auth.casdoor.issuer"))
	clientID := strings.TrimSpace(viper.GetString("auth.casdoor.client_id"))
	if issuer == "" || clientID == "" {
		logger.Warn("Casdoor OIDC 参数不完整，跳过 OIDC verifier 初始化")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		logger.Errorf("Casdoor OIDC provider 初始化失败: %v", err)
		return
	}

	OIDCVerifier = provider.Verifier(&oidc.Config{ClientID: clientID})
	logger.Info("Casdoor OIDC verifier 初始化完成")
}

func (db *Database) Close() {
	sqlDB, err := db.Self.DB()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = sqlDB.Close()
	//err := DB.Self.Close()
	if err != nil {
		logger.Error("关闭连接失败，错误信息: %s", err)
	}

	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}
