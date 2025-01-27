package cmd

import (
	"fmt"

	"github.com/labstack/echo/v4"

	_ "github.com/r4f3t/messagesender/docs"
	"github.com/r4f3t/messagesender/helper"
	"github.com/r4f3t/messagesender/internal/message"
	"github.com/r4f3t/messagesender/internal/message/controller"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type api struct {
	instance *echo.Echo
	command  *cobra.Command
	Port     string
}

// apiCmd represents the api command
var apiCmd = &api{
	command: &cobra.Command{
		Use:   "api",
		Short: "",
		Long:  "",
	},
	Port: "8080",
}

func init() {
	RootCommand.AddCommand(apiCmd.command)
	apiCmd.command.Flags().StringVarP(&apiCmd.Port, "port", "p", "8080", "Service Port")
	apiCmd.instance = echo.New()

	apiCmd.command.RunE = func(cmd *cobra.Command, args []string) error {
		db := helper.InitializeDatabase()
		redisDb := helper.InitializeRedis()

		apiCmd.instance.GET("/swagger/*", echoSwagger.WrapHandler)

		messageRepository := message.NewRepository(db)
		messageService := message.NewService(messageRepository, redisDb)
		messageService.StartProcessingMessages()

		controller.MakeHandler(apiCmd.instance, controller.NewController(messageService))

		apiCmd.instance.Logger.Fatal(apiCmd.instance.Start(fmt.Sprintf(":%s", apiCmd.Port)))
		return nil
	}

}
