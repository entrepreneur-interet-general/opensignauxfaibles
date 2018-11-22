package main

import (
	"github.com/chrnin/ganboard"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Kanboard -> ganboard.Client in gingonic context
func Kanboard() gin.HandlerFunc {

	endpoint := viper.GetString("KANBOARD_ENDPOINT")
	username := viper.GetString("KANBOARD_USERNAME")
	password := viper.GetString("KANBOARD_PASSWORD")

	client := ganboard.Client{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}

	projects, err := client.GetAllProjects()
	// logerror.Println("Connexion à Kanboard et récupération des projets")

	if err != nil {
		// log.Panic("Kanboard n'est pas disponible, avortement du démarrage")
	}

	for _, project := range projects {
		if project.Identifier == "SIGNAUXFAIBLES" {
			client.Project = project
		}
	}

	if client.Project.ID == 0 {
		projectID, err := client.CreateProject(ganboard.ProjectParams{
			Name:        "Signaux-Faibles™",
			Identifier:  "signauxfaibles",
			Description: "Test project for SignauxFaibles",
			OwnerID:     1,
		})

		if err != nil {
			// log.Panic("Kanboard n'est pas disponible, avortement du démarrage")
		}
		client.Project, err = client.GetProjectByID(projectID)

		if err != nil {
			// log.Panic("Kanboard n'est pas disponible, avortement du démarrage")
		}
	}

	return func(c *gin.Context) {
		c.Set("KBCLIENT", &client)
		c.Next()
	}
}

func getKBTasks(c *gin.Context) {
	kb := c.Keys["KBCLIENT"].(*ganboard.Client)
	allTasks, err := kb.GetAllTasks(1, 1)
	if err == nil {
		c.JSON(200, allTasks)
	} else {
		spew.Dump(err)
		c.JSON(500, err)
	}
}

func listProjects(c *gin.Context) {
	kb := c.Keys["KBCLIENT"].(*ganboard.Client)

	project, _ := kb.GetAllProjects()
	c.JSON(200, project)
}
