package internal

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func SetLogLevel(level string) {
	switch level {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Warnln("Unrecognized minimum log level; using 'info' as default")
		log.SetLevel(log.InfoLevel)
	}
}

func SetLogFormatter(timestampFormat string) {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = timestampFormat
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Md5sumFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Md5sum(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// FindMatches finds regular expression matches in a key/value based
// text (ini files, for example), and returns a map with them.
// If the matched key has spaces, they will be replaced with underscores
// If the same keys is found multiple times, the entry of the map will
// have a list as value with all of the matched values
// The pattern must have 2 groups. For example: `(.+)=(.*)`
func FindMatches(pattern string, text []byte) map[string]interface{} {
	configMap := make(map[string]interface{})

	r := regexp.MustCompile(pattern)
	values := r.FindAllStringSubmatch(string(text), -1)
	for _, match := range values {
		key := strings.Replace(match[1], " ", "_", -1)
		if _, ok := configMap[key]; ok {
			switch configMap[key].(type) {
			case string:
				configMap[key] = []interface{}{configMap[key]}
			}
			configMap[key] = append(configMap[key].([]interface{}), match[2])
		} else {
			configMap[key] = match[2]
		}
	}
	return configMap
}

func CRC32hash(input []byte) int {
	crc32Table := crc32.MakeTable(crc32.IEEE)
	return int(crc32.Checksum(input, crc32Table))

}

// Repeat executes a function at a given interval.
// the first tick runs immediately
func Repeat(operation string, tick func(), interval time.Duration, ctx context.Context) {
	tick()

	ticker := time.NewTicker(interval)
	msg := fmt.Sprintf("Next execution for operation %s in %s", operation, interval)
	log.Debugf(msg)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			tick()
			log.Debugf(msg)
		case <-ctx.Done():
			return
		}
	}
}
