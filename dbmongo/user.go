package main

import (
	"bytes"
	"crypto/rand"
	"dbmongo/lib/engine"
	"errors"
	"fmt"
	"html/template"
	"math/big"
	"net/smtp"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// login
type login struct {
	Email        string `form:"email" json:"email" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	BrowserToken string `form:"browserToken" json:"browserToken"`
	CheckCode    string `form:"checkCode" json:"checkCode"`
}

// AdminUser object utilisateur mongodb
type AdminUser struct {
	ID             engine.AdminID `json:"_id" bson:"_id"`
	HashedPassword []byte         `json:"hashedPassword,omitempty" bson:"hashedPassword,omitempty"`
	HashedRecovery []byte         `json:"hashedRecovery,omitempty" bson:"hashedRecovery,omitempty"`
	TimeRecovery   time.Time      `json:"timeRecovery" bson:"timeRecovery"`
	HashedCode     []byte         `json:"hashedCode,omitempty" bson:"hashedCode,omitempty"`
	TimeCode       time.Time      `json:"timeCode,omitempty" bson:"timeCode,omitempty"`
	Cookies        []string       `json:"cookies" bson:"cookies"`
	Level          string         `json:"level" bson:"level"`
	FirstName      string         `json:"firstName" bson:"firstName"`
	LastName       string         `json:"lastName" bson:"lastName"`
	BrowserTokens  []string       `json:"browserTokens" bson:"browserTokens"`
	Regions        []string       `json:"regions" bson:"regions"`
}

func (user AdminUser) save() error {
	err := engine.Db.DBStatus.C("Admin").Update(bson.M{"_id": user.ID}, user)
	return err
}

// type AdminLevel string
const levelAdmin = "admin"
const levelPowerUser = "powerUser"
const levelUser = "user"

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	email := claims["id"].(string)
	user, err := loadUser(email)
	if err != nil {
		c.JSON(500, "Erreur d'identification")
	}
	return &user
}

func loginUser(username string, password string, browserToken string) (AdminUser, error) {
	var user AdminUser
	if err := engine.Db.DBStatus.C("Admin").Find(bson.M{"_id.type": "credential", "_id.key": username}).One(&user); err != nil {
		return AdminUser{}, err
	}
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))

	_, errToken := readBrowserToken(browserToken)
	if err == nil && errToken == nil {
		return user, nil
	}
	return AdminUser{}, errors.New("nop")
}

func loginUserWithCredentials(username string, password string) (AdminUser, error) {
	var user AdminUser
	if err := engine.Db.DBStatus.C("Admin").Find(bson.M{"_id.type": "credential", "_id.key": username}).One(&user); err != nil {
		return AdminUser{}, err
	}
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == nil {
		return user, nil
	}
	return AdminUser{}, err
}

func loadUser(email string) (AdminUser, error) {
	var user AdminUser
	if err := engine.Db.DBStatus.C("Admin").Find(bson.M{"_id.type": "credential", "_id.key": email}).One(&user); err != nil {
		return AdminUser{}, err
	}
	return user, nil
}

func authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login

	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	password := loginVals.Password
	browserToken := loginVals.BrowserToken
	user, err := loginUser(email, password, browserToken)

	if err == nil {
		return user, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*AdminUser); ok && v.Level == "admin" {
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func unauthorizedHandler(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(AdminUser); ok {
		return jwt.MapClaims{
			"id": v.ID.Key,
		}
	}
	return jwt.MapClaims{}
}

func getCode() int {
	i, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return int(i.Int64())
}

//
// @summary Récupération du mot de passe
// @description Dans le cas d'un oubli du mot de passe, on peut fixer un nouveau mot de passe à partir d'un navigateur identifié
// @description Il faut dans ce cas pouvoir recevoir un code de vérification par mail
// @description Voir /login/recovery/get
// @Tags Authentification
// @accept  json
// @produce  json
// @Params email query string true "Adresse e-mail"
// @Params code query string true "Code de vérification"
// @Params password query string true "Nouveau mot de passe"
// @Params browserToken query string true "Jeton du navigateur"
// @Success 200 {string} string "ok"
// @Router /login/recovery/get [post]
func getRecoveryEmailHandler(c *gin.Context) {
	var request struct {
		Email        string `json:"email"`
		BrowserToken string `json:"browserToken"`
	}
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(400, "Bad Parameters 1")
		return
	}

	email := request.Email
	browser, err := readBrowserToken(request.BrowserToken)
	if err != nil || browser.Email != email {
		c.JSON(400, "Bad Parameters 2")
		return
	}

	err = sendRecoveryEmail(email)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, nil)
	}
}

func sendRecoveryEmail(email string) error {
	user, err := loadUser(email)
	if err == nil {
		recoveryCode := fmt.Sprintf("%06d", getCode())
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(recoveryCode), bcrypt.DefaultCost)
		if err == nil {
			user.HashedRecovery = hashedPassword
			user.TimeRecovery = time.Now()
			err = user.save()

			if err == nil {
				fmt.Println(recoveryCode)

				smtpAddress := viper.GetString("smtpAddress")
				//smtpUser := viper.GetString("smtpUser")
				//smtpPassword := viper.GetString("smtpPass")

				c, err := smtp.Dial(smtpAddress)
				if err != nil {
					spew.Dump(err)
				}
				defer c.Close()
				// Set the sender and recipient.
				c.Mail("do.not.reply@signaux.faibles.fr")
				c.Rcpt(email)

				// Send the email body.
				wc, err := c.Data()
				if err != nil {
					spew.Dump(err)
				}
				defer wc.Close()
				buf := bytes.NewBufferString(`Bonjour,
				
				suite à votre demande de récupération de mot de passe sur l'applicatif Signaux Faibles, voici votre code de vérification:` + recoveryCode + `
				
				Cordialement,
			
				l'équipe Signaux-Faibles.
			
				ps: si vous n'êtes pas à l'origine de cette tentative, nous vous prions d'en faire part à l'adresse contact@signaux-faibles.beta.gouv.fr
				
				`)

				if _, err = buf.WriteTo(wc); err != nil {
					spew.Dump(err)
				}
				return err

			}
		}
	} else {
		fmt.Println("error: " + err.Error())
	}

	return nil
}

func getRegions() map[string]string {
	regions := map[string]string{
		"ARA": "Auvergne-Rhône-Alpes",
		"BFC": "Bourgogne-Franche-Comté",
		"BRE": "Bretagne",
		"CVL": "Centre-Val de Loire",
		"COR": "Corse",
		"GES": "Grand Est",
		"HDF": "Hauts-de-France",
		"IDF": "Île-de-France",
		"NOR": "Normandie",
		"NAQ": "Nouvelle-Aquitaine",
		"OCC": "Occitanie",
		"PDL": "Pays de la Loire",
		"PAC": "Provence-Alpes-Côte d'Azur",
	}
	return regions
}

//
// @summary Récupération du mot de passe
// @description Dans le cas d'un oubli du mot de passe, on peut fixer un nouveau mot de passe à partir d'un navigateur identifié
// @description L'utilisation de ce service permet de vérifier le code envoyé à l'utilisateur sur sa messagerie
// @description Ce code est validé en utilisant le hash (bcrypt) stocké en base par le service /login/recovery/get
// @Tags Authentification
// @accept  json
// @produce  json
// @Params email query string true "Adresse e-mail"
// @Params code query string true "Code de vérification"
// @Params password query string true "Nouveau mot de passe"
// @Params browserToken query string true "Jeton du navigateur"
// @Success 200 {string} string "Le serveur a accepté la vérification, le mot de passe est changé"
// @Failure 400 {string} string "Les paramètres sont incorrects"
// @Failure 500 {string} string "La vérification du code a échoué (raison non spécifiée)"
// @Router /login/recovery/setPassword [post]
func checkRecoverySetPassword(c *gin.Context) {
	var request struct {
		Email        string `json:"email"`
		RecoveryCode string `json:"code"`
		Password     string `json:"password"`
		BrowserToken string `json:"browserToken"`
	}
	err := c.ShouldBind(&request)

	if err != nil {
		c.JSON(400, "Bad Parameters")
	}

	browser, err := readBrowserToken(request.BrowserToken)
	if err != nil {
		c.JSON(400, "Bad Parameters")
	}

	email := request.Email
	code := request.RecoveryCode
	password := request.Password
	if browser.Email != email {
		c.JSON(400, "Bad Parameters")
	}

	user, err := loadUser(email)
	if err != nil {
		c.JSON(400, "Bad Parameters")
	}
	err = bcrypt.CompareHashAndPassword(user.HashedRecovery, []byte(code))
	if err != nil {
		c.JSON(500, "Server side error")
	}
	user.HashedRecovery = nil
	user.TimeRecovery = time.Time{}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(500, "Server side error")
	}
	user.HashedPassword = hashedPassword
	err = user.save()

	if err != nil {
		c.JSON(500, "Server side error")
	}
}

//
// @summary Envoie un code de vérification par EMail
// @description Le code n'est envoyé que si le mot de passe est valide.
// @description Le cas échéant, un email avertissant d'une tentative est envoyé.
// @description Pour éviter les tentatives de forçage de mot de passe, ce service ne renvoie jamais d'échec.
// @Tags Authentification
// @accept  json
// @produce  json
// @Param email query string true "Adresse EMail"
// @Param password query string true "Mot de passe"
// @Success 200 {string} string "ok"
// @Router /login/get [post]
func sendMail(sender string, recipient string, title string, body string) error {
	smtpAddress := viper.GetString("smtpAddress")
	smtpConnection, err := smtp.Dial(smtpAddress)
	if err != nil {
		return err
	}
	defer smtpConnection.Close()
	smtpConnection.Mail(sender)
	smtpConnection.Rcpt(recipient)

	emailObject, err := smtpConnection.Data()
	if err != nil {
		return err
	}
	defer emailObject.Close()

	return nil
}

func loginGetHandler(c *gin.Context) {
	var loginVals login

	if err := c.ShouldBind(&loginVals); err != nil {
		c.JSON(401, "requête malformée")
		return
	}

	loginGet(loginVals)

	c.JSON(200, "ok")
}

func loginGet(login login) error {
	email := login.Email
	password := login.Password
	user, err := loginUserWithCredentials(email, password)

	mailTemplate, _ := template.New("loginMail").Parse(`
	Subject: Signaux-Faibles – votre code de vérification
	Content-Type: text/plain; charset=us-ascii; format=flowed
	Content-Transfer-Encoding: 7bit
	
	Bonjour,
	suite à votre tentative d'identification sur l'applicatif Signaux Faibles, voici votre code de vérification:
	{{.CheckCode}}
	
	Cordialement,
	l'équipe Signaux-Faibles.
				
	ps: si vous n'êtes pas à l'origine de cette tentative, nous vous prions d'en faire part à l'adresse contact@signaux-faibles.beta.gouv.fr`)

	if err == nil {
		code := struct{ checkCode string }{}
		code.checkCode = fmt.Sprintf("%06d", getCode())
		hashedCode, err := bcrypt.GenerateFromPassword([]byte(code.checkCode), bcrypt.DefaultCost)
		if err == nil {
			user.HashedCode = hashedCode
			user.TimeCode = time.Now()
			err = user.save()
			if err == nil {
				fmt.Println(code.checkCode)

				smtpAddress := viper.GetString("smtpAddress")
				//smtpUser := viper.GetString("smtpUser")
				//smtpPassword := viper.GetString("smtpPass")

				c, err := smtp.Dial(smtpAddress)
				if err != nil {
					return err
				}
				defer c.Close()

				// Set the sender and recipient.
				c.Mail("Signaux Faibles <do.not.reply@signaux.faibles.fr>")
				c.Rcpt(email)

				// Send the email body.
				wc, err := c.Data()
				if err != nil {
					spew.Dump(err)
				}

				mailTemplate.Execute(wc, code)
				wc.Close()
				return err
			}
		}
	} else {
		fmt.Println("error: " + err.Error())
	}
	return err
}

//
// @summary Vérification du code temporaire renvoyé à l'utilisateur
// @description Fournit en retour un jeton de navigateur
// @Tags Authentification
// @accept  json
// @produce  json
// @Param email query string true "Adresse EMail"
// @Param password query string true "Mot de passe"
// @Param checkCode query string true "Code de vérification"
// @Success 200 {string} string "Le mot de passe et le code de vérification correspondent aux valeurs hashées en base"
// @Router /login/check [post]
func loginCheckHandler(c *gin.Context) {
	var loginVals login
	c.ShouldBind(&loginVals)

	email := loginVals.Email
	password := loginVals.Password
	checkCode := loginVals.CheckCode

	err := loginCheck(email, password, checkCode)

	if err != nil {
		c.JSON(401, "Erreur d'authentification")
	} else {
		browser := Browser{
			IP:      c.ClientIP(),
			Created: time.Now(),
			Email:   email,
			// TODO: identifier les navigateurs
			Name: "",
		}
		browserToken, _ := forgeBrowserToken(browser)
		c.JSON(200, browserToken)
	}
}

func loginCheck(email string, password string, checkCode string) error {
	user, err := loginUserWithCredentials(email, password)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedCode, []byte(checkCode))
	if err != nil {
		return err
	}

	user.HashedCode = nil
	user.TimeCode = time.Time{}
	err = user.save()
	return err
}

//
// @summary Rafraichir le jeton d'identification
// @description Nécessite 3 informations: email, mot de passe et jeton de navigateur
// @description Le jeton de navigateur n'a pas de limite de validité et peut être conservé
// @Tags Authentification
// @accept  application/json
// @produce  application/json
// @Param email query string true "Adresse Email" "gnigni"
// @Param password query string true "Mot de Passe"
// @Param browserToken query string true "Token navigateur"
// @Success 200 {string} string "{<br/>'code': 200,<br/>'expire':'2019-01-21T11:16:15+01:00',<br/>'token':'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…'<br/>}"
// @Failure 401 {string} string "{<br/>'code': 401,<br/>'message': 'incorrect Username or Password'<br/>}"
// @Router /login [post]
func dummyLogin() {}

//
// @summary Obtenir un jeton d'identification
// @description Fournit un jeton avec nouvelle date de validité en échange d'un jeton encore valide
// @Tags Authentification
// @accept  application/json
// @produce  application/json
// @Security ApiKeyAuth
// @Success 200 {string} string "{<br/>'code': 200,<br/>   'expire': '2019-01-21T12:16:50+01:00',<br/>'token': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…'<br/>}"
// @Failure 401 {string} string "{<br/>'code': 401,<br/>'message': 'cookie token is empty'<br/>}"
// @Router /api/refreshToken [get]
func dummyRefreshtoken() {}
