package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nickemma/go-jwt/initializers"
	"github.com/nickemma/go-jwt/models"
)

func Auth(c *gin.Context) {
	 // get the cookie token off the req body

		tokenString, err := c.Cookie("Authorization")

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// decode and validate the token

		 token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			
	  // Don't forget to validate the alg is what you expect:
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		   return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	   }

	 // Get Secret in a []byte containing your secret
	   return []byte(os.Getenv("SECRET")), nil
    })
				
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
						// check for expiration
						if float64(time.Now().Unix()) > claims["exp"].(float64) {
	     			c.AbortWithStatus(http.StatusUnauthorized)
						}

		 // find the user with the token sub decoded

						var user models.User
						initializers.DB.First(&user, claims["sub"])

						if user.ID == 0 {	     			
							c.AbortWithStatus(http.StatusUnauthorized)
						}

		 // attach the decoded req 
						c.Set("user", user)

		 // continue 
		  c.Next()
				
    } else {
	     			c.AbortWithStatus(http.StatusUnauthorized)
    }
	
}