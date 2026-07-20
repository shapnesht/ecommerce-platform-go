package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLmode  string
}

type JWTConfig struct {
	Secret             string
	TokenExpiry        time.Duration
	RefreshTokenExpiry time.Duration
}

type AWSConfig struct {
	Region       string
	AccessKeyID  string
	SecretKey    string
	S3BucketName string
	S3Endpoint   string
}

type UploadConfig struct {
	UploadPath string
	MaxSize    int64
	// Cab be s3 or local
	UploadProvider string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	jwtExpiry, _ := time.ParseDuration(GetEnv("JWT_EXPIRY", "24h"))
	refreshTokenExpiry, _ := time.ParseDuration(GetEnv("REFRESH_TOKEN_EXPIRY", "72h"))
	maxUploadSize, _ := strconv.ParseInt(GetEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	return &Config{
		Server: ServerConfig{
			Port:    GetEnv("PORT", "8080"),
			GinMode: GetEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			Username: GetEnv("DB_USERNAME", "postgres"),
			Password: GetEnv("DB_PASSWORD", "password"),
			Name:     GetEnv("DB_NAME", "ecommerce_shop"),
			SSLmode:  GetEnv("DB_SSLmode", "disable"),
		},
		JWT: JWTConfig{
			TokenExpiry:        jwtExpiry,
			RefreshTokenExpiry: refreshTokenExpiry,
			Secret:             GetEnv("JWT_SECRET_KEY", "secretkey"),
		},
		AWS: AWSConfig{
			Region:       GetEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:  GetEnv("AWS_ACCESS_KEY_ID", "test"),
			SecretKey:    GetEnv("AWS_SECRET_ACCESS_KEY", "secret"),
			S3BucketName: GetEnv("AWS_S3_BUCKET_NAME", "ecommerce-uploads"),
			S3Endpoint:   GetEnv("AWS_S3_ENDPOINT", "localhost:9000"),
		},
		Upload: UploadConfig{
			MaxSize:        maxUploadSize,
			UploadPath:     GetEnv("UPLOAD_PATH", "./uploads"),
			UploadProvider: GetEnv("UPLOAD_PROVIDER", "s3"),
		},
	}, nil

}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
