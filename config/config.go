package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	GoroutinesCount int
	URLs            []string
}

const (
	httpPrefix  = "http://"
	httpsPrefix = "https://"
)

func ReadConfig() (Config, error) {
	count := 0
	flag.IntVar(&count, "parallel", 10, "Count of workers that will send requests")
	flag.Parse()

	return NewConfig(count, flag.Args())
}

func NewConfig(goroutinesCount int, urls []string) (Config, error) {
	if goroutinesCount <= 0 {
		return Config{}, errors.New("count should be greater than 0")
	}

	if len(urls) == 0 {
		return Config{}, errors.New("no urls set")
	}

	result := Config{
		GoroutinesCount: goroutinesCount,
		URLs:            make([]string, 0, len(urls)),
	}
	for _, url := range urls {
		if strings.HasPrefix(url, httpPrefix) || strings.HasPrefix(url, httpsPrefix) {
			result.URLs = append(result.URLs, url)
			continue
		}
		result.URLs = append(result.URLs, fmt.Sprintf("%s%s", httpPrefix, url))
	}

	return result, nil
}
