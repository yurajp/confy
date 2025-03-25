package confy

import (
  	"os"
  	"fmt"
  	"sort"
    "bufio"
	"strings"
)

func Exist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func mapping(path string) (map[string]string, error) {
	mp := make(map[string]string)
	ex := Exist(path)
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0640)
	if err != nil {
		return mp, err
	}
	if !ex {
		return mp, nil
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(ln, "#") {
			continue
		}
		vars := strings.SplitN(ln, "=", 2)
		mp[vars[0]] = vars[1]
	}
  if err := sc.Err(); err != nil {
        return map[string]string{}, err
    }
  return mp, nil
}

func LoadEnv(path string) error {
	mp, err := mapping(path)
	if err != nil {
		return err
	}
	for k, val := range mp {
		os.Setenv(k, val)
	}
	return nil
}

func AddVar(path, key, val string) error {
	mp, err := mapping(path)
  f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0644)
  defer f.Close()
  mp[key] = val
  keys := []string{}
  for k, _ := range mp {
	  keys = append(keys, k)
  }
  sort.Strings(keys)
  for _, k := range keys {
	  ln := fmt.Sprintf("%s=%s\n", k, mp[k])
	  _, err = f.WriteString(ln)
	  if err != nil {
		  return err
	  }
  }
  return nil
}


