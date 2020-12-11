package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"log"
	"net/url"
	"strings"

	"net/http"
	"os"
	"strconv"
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
	AccessUuid string
	UserId     uint64
}

const AccessTokenKey string = "access_token"
const RefreshTokenKey string = "refresh_token"

func IsRequestAuthorized(request *http.Request) (bool, error) {
	tokenAuth, err := ExtractTokenMetadata(request)
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
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		return false, err
	}
	refreshToken := mapToken[RefreshTokenKey]

	//verify the token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
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
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return false, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return false, err
		}
		//Delete the previous Refresh Token
		deleted, delErr := DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			return false, err
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(userId)
		if createErr != nil {
			return false, err
		}
		//save the tokens metadata to redis
		saveErr := CreateAuth(userId, ts)
		if saveErr != nil {
			return false, err
		}
		tokens := map[string]string{
			AccessTokenKey:  ts.AccessToken,
			RefreshTokenKey: ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     AccessTokenKey,
			Value:    url.QueryEscape(ts.AccessToken),
			HttpOnly: true,
			Secure:   false,
		})

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     RefreshTokenKey,
			Value:    url.QueryEscape(ts.RefreshToken),
			HttpOnly: true,
			Secure:   false,
		})

		return true, nil
	} else {
		return false, nil
	}
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := RedisClient.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func FetchAuth(authD *accessDetails) (uint64, error) {
	userid, err := RedisClient.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func extractTokensFromCookie(r *http.Request) map[string]string {
	result := make(map[string]string, 2)
	accessTokenCookie, err := r.Cookie(AccessTokenKey)
	if err != nil {
		log.Panic(err)
		return result
	}
	refreshTokenCookie, err := r.Cookie(RefreshTokenKey)
	if err != nil {
		log.Panic(err)
		return result
	}

	result[AccessTokenKey] = accessTokenCookie.Value
	result[RefreshTokenKey] = refreshTokenCookie.Value
	return result
}

func ExtractAccessTokenFromCookie(r *http.Request) string {
	cookies := extractTokensFromCookie(r)
	return cookies[AccessTokenKey]
}

func ExtractTokenFromBearToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractAccessTokenFromCookie(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*accessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &accessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := RedisClient.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := RedisClient.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func CreateToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}
