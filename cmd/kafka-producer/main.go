package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	produceTools "gitlab.tubecorporate.com/push/kafka-producer/internal/produce-tools"

	server "gitlab.tubecorporate.com/push/kafka-producer/internal/server"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Too little arguments, please pass path to config file")
		os.Exit(1)
	}

	Port := getPort()

	cfg := getConfigMap(os.Args[1])

	producer := createProducer(cfg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig.String(), "called. Closing server and producer...")
			producer.Close()
			os.Exit(1)
		}
	}()

	fmt.Println("Starting server on port", Port)
	if err := server.LaunchFastHTTPServer(Port, producer); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getPort() string {
	Port := ":8080"
	if port := os.Getenv("Port"); len(port) > 2 && port[0] == ':' {
		Port = port
		fmt.Println("Passed custom", Port, "port")
	}
	return Port
}

func getConfigMap(filePath string) map[string]interface{} {
	cfg, err := loadConfiguration(filePath)
	if err != nil {
		fmt.Println("Error parse producer config file", err.Error())
		os.Exit(1)
	}
	return cfg
}

func createProducer(cfg map[string]interface{}) *produceTools.Producer {
	Producer, err := produceTools.NewKafkaProducer(cfg)
	if err != nil {
		fmt.Println("Error create producer", err.Error())
		os.Exit(2)
	}
	fmt.Println("Producer created")
	return Producer
}

func loadConfiguration(file string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)

	return result, err
}
