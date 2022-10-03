package helpers

import (
	"context"
	"errors"
	"fmt"
	"nft-raffle/database"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenHelper interface {
	GenerateAllTokens(email, firstName, lastName, uid, userRole string, is_email_verified bool) (signedToken, signedRefreshToken string, err error)
	UpdateAllTokens(signedToken, signedRefreshToken, userId string) error
	ValidateAccessToken(signedToken string) (claims *SignedDetails, err error)
	ValidateRefreshToken(signedToken string) (claims *SignedDetails, err error)
	SetBlacklistAccessTokenUserId(userId string) error
	SetBlacklistRefreshTokenUserId(userId string) error
	SetBlacklistAccessAndRefreshTokenUserId(userId string) error
	GetBlacklistAccessTokenUserId(userId string) (int64, error)
	GetBlacklistRefreshTokenUserId(userId string) (int64, error)
}

type tokenHelperStruct struct{}
type SignedDetails struct {
	Email             string
	First_name        string
	Last_name         string
	Uid               string
	User_role         string
	Is_email_verified bool
	jwt.StandardClaims
}

const (
	blacklistAccessToken  string = "blacklist_access_token"
	blacklistRefreshToken string = "blacklist_refresh_token"
)

var (
	nftRaffleDb       database.NftRaffleMongoDbConnection = database.NewNftRaffleMongoDbConnection()
	nftRaffleDbClient *mongo.Client                       = nftRaffleDb.DBClient()
	userCollection    *mongo.Collection                   = nftRaffleDb.OpenCollection(nftRaffleDbClient, "user")

	redisDb     database.RedisConnection = database.NewRedisConenction()
	redisClient                          = redisDb.RedisClient()

	accessTokenSecretKey  = dotEnvHelperImpl.GetEnvVariable("MY_ACCESS_TOKEN_SECRET_KEY")
	refreshTokenSecretKey = dotEnvHelperImpl.GetEnvVariable("MY_REFRESH_TOKEN_SECRET_KEY")
	accessTokenTTL        = dotEnvHelperImpl.GetEnvVariable("ACCESS_TOKEN_TTL")
	refreshTokenTTL       = dotEnvHelperImpl.GetEnvVariable("REFRESH_TOKEN_TTL")
)

func NewTokenHelper() TokenHelper {
	return &tokenHelperStruct{}
}

func (t *tokenHelperStruct) GenerateAllTokens(email, firstName, lastName, uid, userRole string, is_email_verified bool) (signedToken, signedRefreshToken string, err error) {
	accessTokenTTLHoursInt, err := strconv.Atoi(accessTokenTTL)

	if err != nil {
		return "", "", err
	}

	refreshTokenTTLHoursInt, err := strconv.Atoi(refreshTokenTTL)

	if err != nil {
		return "", "", err
	}

	claims := &SignedDetails{
		Email:             email,
		First_name:        firstName,
		Last_name:         lastName,
		Uid:               uid,
		User_role:         userRole,
		Is_email_verified: is_email_verified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(accessTokenTTLHoursInt)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
			Subject:   uid,
		},
	}

	refreshClaims := &SignedDetails{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(refreshTokenTTLHoursInt)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
			Subject:   uid,
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(accessTokenSecretKey))

	if err != nil {
		return
	}

	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(refreshTokenSecretKey))

	if err != nil {
		return
	}

	return signedToken, signedRefreshToken, err
}

func (t *tokenHelperStruct) UpdateAllTokens(signedToken, signedRefreshToken, userId string) error {
	var updateObj bson.D

	updateObj = append(updateObj, bson.E{Key: "access_token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	if err != nil {
		return err
	}

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	_, updateError := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()

	if updateError != nil {
		return updateError
	}

	return nil
}

func (t *tokenHelperStruct) ValidateAccessToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(accessTokenSecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		return nil, errors.New("invalid access token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("access token has expired")
	}

	return claims, nil
}

func (t *tokenHelperStruct) ValidateRefreshToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(refreshTokenSecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		return nil, errors.New("invalid refresh token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}

func (t *tokenHelperStruct) SetBlacklistAccessTokenUserId(userId string) error {
	accessTokenTTLHoursInt, err := strconv.Atoi(accessTokenTTL)

	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s:%s", blacklistAccessToken, "user_id", userId)
	val := time.Now().Local().Add(time.Hour * time.Duration(accessTokenTTLHoursInt)).Unix()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	err = redisClient.Set(ctx, key, val, time.Hour*time.Duration(accessTokenTTLHoursInt)).Err()
	defer cancel()

	if err != nil {
		return err
	}

	return nil
}

func (t *tokenHelperStruct) SetBlacklistRefreshTokenUserId(userId string) error {
	refreshTokenTTLHoursInt, err := strconv.Atoi(refreshTokenTTL)

	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s:%s", blacklistRefreshToken, "user_id", userId)
	val := time.Now().Local().Add(time.Hour * time.Duration(refreshTokenTTLHoursInt)).Unix()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	err = redisClient.Set(ctx, key, val, time.Hour*time.Duration(refreshTokenTTLHoursInt)).Err()
	defer cancel()

	if err != nil {
		return err
	}

	return nil
}

func (t *tokenHelperStruct) SetBlacklistAccessAndRefreshTokenUserId(userId string) error {
	accessTokenTTLHoursInt, err := strconv.Atoi(accessTokenTTL)

	if err != nil {
		return err
	}

	refreshTokenTTLHoursInt, err := strconv.Atoi(refreshTokenTTL)

	if err != nil {
		return err
	}

	blacklistAccessTokenKey := fmt.Sprintf("%s:%s:%s", blacklistAccessToken, "user_id", userId)
	blacklistAccessTokenVal := time.Now().Local().Add(time.Hour * time.Duration(accessTokenTTLHoursInt)).Unix()

	blacklistRefreshTokenKey := fmt.Sprintf("%s:%s:%s", blacklistRefreshToken, "user_id", userId)
	blacklistRefreshTokenVal := time.Now().Local().Add(time.Hour * time.Duration(refreshTokenTTLHoursInt)).Unix()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	pipe := redisClient.TxPipeline()

	pipe.Set(ctx, blacklistAccessTokenKey, blacklistAccessTokenVal, time.Hour*time.Duration(accessTokenTTLHoursInt))
	defer cancel()

	pipe.Set(ctx, blacklistRefreshTokenKey, blacklistRefreshTokenVal, time.Hour*time.Duration(refreshTokenTTLHoursInt))
	defer cancel()

	_, err = pipe.Exec(ctx)
	defer cancel()

	if err != nil {
		return err
	}

	return nil
}

func (t *tokenHelperStruct) GetBlacklistAccessTokenUserId(userId string) (int64, error) {
	key := fmt.Sprintf("%s:%s:%s", blacklistAccessToken, "user_id", userId)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	val, err := redisClient.Get(ctx, key).Result()
	defer cancel()

	if err == redis.Nil {
		return -1, nil
	} else if err != nil {
		return -1, err
	}

	unixTime, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		return -1, err
	}

	return unixTime, nil
}

func (t *tokenHelperStruct) GetBlacklistRefreshTokenUserId(userId string) (int64, error) {
	key := fmt.Sprintf("%s:%s:%s", blacklistRefreshToken, "user_id", userId)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	val, err := redisClient.Get(ctx, key).Result()
	defer cancel()

	if err == redis.Nil {
		return -1, nil
	} else if err != nil {
		return -1, err
	}

	unixTime, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		return -1, err
	}

	return unixTime, nil
}
