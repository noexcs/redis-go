package config

import (
	"bufio"
	"github.com/noexcs/redis-go/log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ServerProperties struct {
	Bind string
	Port int

	Requirepass string

	Debug bool
}

var defaultProperties *ServerProperties
var Properties *ServerProperties
var configFile = "redis.conf"

func init() {
	defaultProperties = &ServerProperties{Bind: "0.0.0.0", Port: 6397}
}

func Setup() {
	if fileExist(configFile) {
		Properties = parseConfigFile(configFile)
	} else {
		Properties = defaultProperties
	}
}

func parseConfigFile(filePath string) *ServerProperties {
	file, _ := os.Open(filePath)
	scanner := bufio.NewScanner(file)

	config := &ServerProperties{}
	userSetMap := make(map[string]string)
	for scanner.Scan() {
		// remove leading and trailing space
		line := strings.TrimSpace(scanner.Text())
		// skip comment which starts with '#'
		if len(line) < 1 || (len(line) >= 1 && line[0] == '#') {
			continue
		}
		splits := strings.Fields(strings.ToLower(line))
		// config syntax is the pattern "key value" strictly
		if len(splits) != 2 {
			continue
		}
		key := splits[0]
		value := splits[1]
		userSetMap[key] = value
	}

	// there is an error spawned during reading the "redis.conf"
	if err := scanner.Err(); err != nil {
		log.FatalWithLocation(err)
	}

	tPtr := reflect.TypeOf(config)
	vPtr := reflect.ValueOf(config)

	numField := tPtr.Elem().NumField()
	for i := 0; i < numField; i++ {
		structureField := tPtr.Elem().Field(i)
		fieldValue := vPtr.Elem().Field(i)

		// ServerProperties' field name --> key (in redis.conf)
		keyName := strings.ToLower(structureField.Name)

		userSetValue, ok := userSetMap[keyName]
		// if user have set the value if key keyName
		if ok {
			switch structureField.Type.Kind() {
			case reflect.String:
				fieldValue.SetString(userSetValue)
			case reflect.Int:
				parseInt, err := strconv.ParseInt(userSetValue, 10, 64)
				if err == nil {
					fieldValue.SetInt(parseInt)
				}
			case reflect.Bool:
				fieldValue.SetBool(userSetValue == "yes")
			}
		}
	}

	return config
}

func fileExist(file string) bool {
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.WithLocation("Config file %s does not exist.", file)
	}

	return err == nil && !fileInfo.IsDir()
}
