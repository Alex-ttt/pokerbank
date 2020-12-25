package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"strings"

	"net/http"
	"os"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type accessDetails struct {
	Uuid      string
	UserLogin string
}

const AccessTokenKey string = "access_token"
const RefreshTokenKey string = "refresh_token"
const AccessSecretKey string = "ACCESS_SECRET"
const RefreshSecretKey string = "REFRESH_SECRET"

func IsRequestAuthorized(request *http.Request) (bool, error) {
	tokenAuth, err := ExtractTokenMetadata(request, true)
	if err != nil {
		return false, err
	}

	_, err = FetchAuth(tokenAuth)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Refresh(c *gin.Context) (bool, error) {
	refreshToken := ExtractRefreshTokenFromCookie(c.Request)

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(RefreshSecretKey)), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return false, err
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false, err
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return false, err
		}
		userLogin := claims["user_login"].(string)

		storedTokenValue, redisErr := GetFromRedis(RedisClient, claims["refresh_uuid"].(string))
		if redisErr != nil || storedTokenValue != userLogin {
			return false, err
		}

		//Delete the previous Refresh Token
		deleted, delErr := DeleteAuth(refreshUuid)

		if delErr != nil || deleted == 0 { //if any goes wrong
			return false, err
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(userLogin)
		if createErr != nil {
			return false, err
		}

		saveErr := CreateAuth(userLogin, ts)
		if saveErr != nil {
			return false, err
		}
		tokens := map[string]string{
			AccessTokenKey:  ts.AccessToken,
			RefreshTokenKey: ts.RefreshToken,
		}

		SetTokensToResponseCookie(&c.Writer, &tokens)

		return true, nil
	} else {
		return false, nil
	}
}

func SetTokensToResponseCookie(writer *gin.ResponseWriter, tokens *map[string]string) {
	http.SetCookie(*writer, &http.Cookie{
		Name:     AccessTokenKey,
		Value:    url.QueryEscape((*tokens)[AccessTokenKey]),
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(time.Minute * 15),
	})
	http.SetCookie(*writer, &http.Cookie{
		Name:     RefreshTokenKey,
		Value:    url.QueryEscape((*tokens)[RefreshTokenKey]),
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().AddDate(0, 0, 90),
	})
}

func DeleteAuth(givenUuid string) (int64, error) {
	err := DeleteFromRedis(RedisClient, givenUuid)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func FetchAuth(authD *accessDetails) (string, error) {
	userLogin, err := GetFromRedis(RedisClient, authD.Uuid)
	if err != nil || len(userLogin) == 0 {
		return "", err
	}

	return userLogin, nil
}

func extractTokensFromCookie(r *http.Request) map[string]string {
	result := make(map[string]string, 2)
	accessTokenCookie, err := r.Cookie(AccessTokenKey)
	if err == nil {
		result[AccessTokenKey] = accessTokenCookie.Value
	}

	refreshTokenCookie, err := r.Cookie(RefreshTokenKey)
	if err == nil {
		result[RefreshTokenKey] = refreshTokenCookie.Value
	}

	return result
}

func ExtractAccessTokenFromCookie(r *http.Request) string {
	cookies := extractTokensFromCookie(r)
	return cookies[AccessTokenKey]
}

func ExtractRefreshTokenFromCookie(r *http.Request) string {
	cookies := extractTokensFromCookie(r)
	return cookies[RefreshTokenKey]
}

//goland:noinspection GoUnusedExportedFunction
func ExtractTokenFromBearToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request, isAccessToken bool) (*jwt.Token, error) {
	var tokenString string
	if isAccessToken {
		tokenString = ExtractAccessTokenFromCookie(r)
	} else {
		tokenString = ExtractRefreshTokenFromCookie(r)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		var secretKey string
		if isAccessToken {
			secretKey = os.Getenv(AccessSecretKey)
		} else {
			secretKey = os.Getenv(RefreshSecretKey)
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//goland:noinspection GoUnusedExportedFunction
func TokenValid(r *http.Request, isAccessToken bool) error {
	token, err := VerifyToken(r, isAccessToken)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request, isAccessToken bool) (*accessDetails, error) {
	token, err := VerifyToken(r, isAccessToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	var tokenUuid string
	if ok && token.Valid {
		if isAccessToken {
			tokenUuid, ok = claims["access_uuid"].(string)
		} else {
			tokenUuid, ok = claims["refresh_uuid"].(string)
		}
		if !ok {
			return nil, err
		}

		userLogin := claims["user_login"].(string)

		return &accessDetails{
			Uuid:      tokenUuid,
			UserLogin: userLogin,
		}, nil
	}
	return nil, err
}

func CreateAuth(userLogin string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := SetRedisKeyWithExpiration(RedisClient, td.AccessUuid, userLogin, at.Sub(now))
	//Set(td.AccessUuid, userLogin, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := SetRedisKeyWithExpiration(RedisClient, td.RefreshUuid, userLogin, rt.Sub(now))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func CreateToken(userLogin string) (*TokenDetails, error) {
	tokenDetails := &TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.AccessUuid = uuid.NewV4().String()

	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 90).Unix()
	tokenDetails.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDetails.AccessUuid
	atClaims["user_login"] = userLogin
	atClaims["exp"] = tokenDetails.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenDetails.AccessToken, err = at.SignedString([]byte(os.Getenv(AccessSecretKey)))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	rtClaims["user_login"] = userLogin
	rtClaims["exp"] = tokenDetails.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenDetails.RefreshToken, err = rt.SignedString([]byte(os.Getenv(RefreshSecretKey)))
	if err != nil {
		return nil, err
	}
	return tokenDetails, nil
}

func EncryptPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func IsPasswordsEqual(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
