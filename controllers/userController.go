package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nickemma/go-jwt/initializers"
	"github.com/nickemma/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
 // Get the email and password from the req body
				 var body struct {
						Email string
						Password string
					}

					if c.Bind(&body) != nil {
							c.JSON(http.StatusBadRequest, gin.H{
										"error": "Failed to read body",
							})
							return
					}

	// hash the password using bcrypt

					hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

					if err != nil {
							c.JSON(http.StatusBadRequest, gin.H{
										"error": "Failed to hash password",
							})
							return
					}

	// create the user

					user := models.User{Email: body.Email, Password: string(hash)}
					result := initializers.DB.Create(&user)

					if result.Error != nil {
							c.JSON(http.StatusBadRequest, gin.H{
										"error": "Failed to create a user",
							})
							return
					}

	// respond to the data created

	  c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {

			// get the email and password of the registered user

				var body struct {
					Email string
					Password string
				}

				if c.Bind(&body) != nil {
						c.JSON(http.StatusBadRequest, gin.H{
										"error": "Invalid credentials",
							})
							return
				}

			// look up for registered user

				var user models.User
				initializers.DB.First(&user, "email = ?", body.Email)

				if user.ID == 0 {
						c.JSON(http.StatusBadRequest, gin.H{
										"error": "Invalid Email or Password",
							})
							return
				}

			// compare the login password with the hashed password

				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
										"error": "Invalid Email or Password",
							})
							return
				}

			// generate a jwt token 

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	     "sub": user.ID,
	     "exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
    })

   // Sign and get the complete encoded token as a string using the secret
   tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
										"error": "Failed to create token",
							})
							return
			}

			// send the response back
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)

			c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	   user, _ := c.Get("user")
				
				c.JSON(http.StatusOK, gin.H{
					"message": user,
				})
}