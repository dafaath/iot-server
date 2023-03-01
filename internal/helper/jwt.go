package helper

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dafaath/iot-server/configs"
	"github.com/dafaath/iot-server/internal/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func SignUserToken(user entities.UserRead) (string, error) {
	config := configs.GetConfig()
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"idUser":   user.IdUser,
		"email":    user.Email,
		"username": user.Username,
		"status":   user.Status,
		"isAdmin":  user.IsAdmin,
		"iat":      time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.JWT.SecretKey))

	return tokenString, err
}

func ValidateUserToken(tokenString string) (user entities.UserRead, err error) {
	config := configs.GetConfig()
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.JWT.SecretKey), nil
	})
	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.IdUser = int(claims["idUser"].(float64))
		user.Email = claims["email"].(string)
		user.Username = claims["username"].(string)
		user.Status = claims["status"].(bool)
		user.IsAdmin = claims["isAdmin"].(bool)
		return user, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return user, fiber.NewError(401, "Token is malformed")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return user, fiber.NewError(401, "Token is expired or not valid yet")
	} else {
		return user, fiber.NewError(401, fmt.Sprintf("Couldn't handle this token: %s", err.Error()))
	}
}

func ValidateUserCredentical(c *fiber.Ctx) (user entities.UserRead, err error) {
	authorizationCookies := c.Cookies("authorization", "")
	authorizationCookies, err = url.QueryUnescape(authorizationCookies)
	if err != nil {
		return user, err
	}

	headers := c.GetReqHeaders()
	authorizationHeaders, haveAuthorizationHeader := headers["Authorization"]

	authorization := ""

	if !haveAuthorizationHeader && authorizationCookies == "" {
		return user, fiber.NewError(401, "Authorization not present")
	}

	if haveAuthorizationHeader {
		authorization = authorizationHeaders
	} else {
		authorization = authorizationCookies
	}

	authorizationSplit := strings.Split(authorization, " ")
	authorizationType := authorizationSplit[0]
	if authorizationType != "Bearer" {
		return user, fiber.NewError(401, "Authorization type is not Bearer, please use 'Bearer {token}' format on your authorization header")
	}

	token := authorizationSplit[1]

	user, err = ValidateUserToken(token)
	if err != nil {
		return user, err
	}

	return user, nil
}
