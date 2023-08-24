package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"gopkg.in/yaml.v2"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type Config struct {
	Pin      string            `yaml:"pin"`
	WavFile  string            `yaml:"wav_file"`
	Pushover map[string]string `yaml:"pushover"`
}

var cfg *Config

func LoadConfig(filename string) (*Config, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func tellPushover() {
	data := url.Values{
		"token":   {cfg.Pushover["token"]},
		"user":    {cfg.Pushover["user"]},
		"message": {cfg.Pushover["message"]},
	}

	resp, err := http.PostForm("https://api.pushover.net:443/1/messages.json", data)
	if err != nil {
		log.Printf("Error sending push notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response when sending push notification: %v", resp.Status)
	}
}

func main() {
	var err error
	cfg, err = LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if _, err := host.Init(); err != nil {
		log.Fatalf("Failed to initialize periph: %v", err)
	}

	p := gpioreg.ByName(cfg.Pin)
	if p == nil {
		log.Fatalf("Failed to find pin: %s", cfg.Pin)
	}

	var count int
	for {
		if p.Read() == gpio.Low {
			count++
		} else {
			count = 0
		}

		if p.Read() == gpio.Low && count == 2 {
			count = 0
			log.Println("Button press")
			tellPushover()
			go func() {
				cmd := exec.Command("aplay", cfg.WavFile)
				if err := cmd.Run(); err != nil {
					log.Printf("Error playing sound: %v", err)
				}
			}()
		}

		time.Sleep(100 * time.Millisecond)

		if count > 0 {
			log.Println(count)
		}
	}
}
