package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/atsushi-matsui/web-authn-example/db"
	"github.com/atsushi-matsui/web-authn-example/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	w *webauthn.WebAuthn
	err error
	userTable *db.UserTable
	sessionTable *db.SessionTable
)

func main() {
	serverAddress := ":8080"

	wConfig := &webauthn.Config{
		RPDisplayName: "Go Webauthn", // Display Name for your site
		RPID: "localhost", // Generally the FQDN for your site
		RPOrigins: []string{"http://localhost" + serverAddress}, // The origin URLs allowed for WebAuthn requests
	}

	if w, err = webauthn.New(wConfig); err != nil {
		fmt.Println(err)
	}

	userTable = db.NewUserTable()
	sessionTable = db.NewSessionTable()

	r := gin.Default()
	r.LoadHTMLGlob("views/*")
	r.GET("/", Index)
	r.GET("/register/begin/:username", BeginRegistration)
	r.POST("/register/finish/:username", FinishRegistration)
	r.GET("/login/begin/:username", BeginLogin)
	r.POST("/login/finish/:username", FinishLogin)
	r.Run(serverAddress)

	log.Println("starting server at", serverAddress)
}

func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func BeginRegistration(c *gin.Context) {
	userName := c.Param("username")
	user, err := userTable.GetUser(userName)

	if err != nil {
		displayName := strings.Split(userName, "@")[0]
		user = domain.NewUser(userName, displayName);
		userTable.PutUser(user)

		log.Println("New User: ", userName, user.GetId())
	}
	
	options, session, err := w.BeginRegistration(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	sessionTable.PutSession(user, session) 

	c.JSON(http.StatusOK, options)
}

func FinishRegistration(c *gin.Context) {
	userName := c.Param("username")
	user, err := userTable.GetUser(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	if err != nil {
		// Handle Error and return.
		log.Println("Parse error: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	sessionData, err := sessionTable.PullSession(user.GetId())
	if err != nil {
		// Handle Error and return.
		log.Println("Not found session: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	credential, err := w.CreateCredential(*user, *sessionData, response)
	if err != nil {
		// Handle Error and return.
		log.Println("CreateCredential error: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user.AddCredential(*credential)

	c.JSON(http.StatusOK, credential)
}

func BeginLogin(c *gin.Context) {
	userName := c.Param("username")
	user, err := userTable.GetUser(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	options, session, err := w.BeginLogin(user)
	if err != nil {
		// Handle Error and return.
		log.Println("BeginLogin error: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	sessionTable.PutSession(user, session)
	
	c.JSON(http.StatusOK, options)
}

func FinishLogin(c *gin.Context) {
	userName := c.Param("username")
	user, err := userTable.GetUser(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	if err != nil {
		// Handle Error and return.
		log.Println("Parse error: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	sessionData, err := sessionTable.PullSession(user.GetId())
	if err != nil {
		// Handle Error and return.
		log.Println("Not found session: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	credential, err := w.ValidateLogin(*user, *sessionData, response)
	if err != nil {
		// Handle Error and return.
		log.Println("Credential validation error: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, credential)
}
