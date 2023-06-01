package main

import (
	"flag"
	"github.com/bingemate/keycloak-service/cmd"
	"github.com/bingemate/keycloak-service/initializers"
	"log"
)

func main() {
	flag.Parse()
	env, err := initializers.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	logFile := initializers.InitLog(env.LogFile)
	defer logFile.Close()
	log.Println("Starting server mode...")
	cmd.Serve(env)
}
