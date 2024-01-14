package common

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/act-gpt/marino/web"

	"github.com/dlclark/regexp2"
	"github.com/joho/godotenv"
	"github.com/lithammer/shortuuid/v4"
)

func DefaultConfig(ptr interface{}, tag string) {
	setDefaults(ptr, tag)
}

func setField(field reflect.Value, defaultVal string) error {

	if !field.CanSet() {
		return fmt.Errorf("Can't set value\n")
	}
	switch field.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	}
	return nil
}

func setDefaults(ptr interface{}, tag string) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("Not a pointer")
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get(tag); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}

		}
	}
	return nil
}

func Files(folder string, include string) []string {
	res := strings.HasSuffix(folder, "/")
	if !res {
		folder = folder + "/"
	}
	entries, err := os.ReadDir(folder)
	var items []string
	if err != nil {
		return items
	}
	for _, e := range entries {
		match, err := regexp.MatchString(include, e.Name())
		if err == nil && match {
			items = append(items, folder+e.Name())
		}
	}
	return items
}

func Contains(slice []string, s string) int {
	for index, value := range slice {
		if value == s {
			return index
		}
	}
	return -1
}

func Filter[T any](s []T, cond func(t T) bool) []T {
	res := []T{}
	for _, v := range s {
		if cond(v) {
			res = append(res, v)
		}
	}
	return res
}

func GetUUID() string {
	return shortuuid.New()
}

func FormatYYYYMM() string {
	now := time.Now().UTC()
	return now.Format("2006-01")
}

func Struct2JSON(it any) map[string]interface{} {
	var config map[string]interface{}
	s, err := json.Marshal(it)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(s, &config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return config
}

type SplitData struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Segment  string `json:"segment"`
}

func addData(m *regexp2.Match) SplitData {
	gps := m.Groups()
	return SplitData{
		Question: strings.TrimSpace(gps[1].Captures[0].String()),
		Answer:   strings.TrimSpace(gps[2].Captures[0].String()),
	}
}

func FormatSplitText(str string) []SplitData {
	var data []SplitData
	reg := regexp2.MustCompile(`(?m)Q\d+:\s*(.*?)\s*A\d+:\s*([\s\S]*?)(?=Q\d+:)`, regexp2.RE2)

	if m, _ := reg.FindStringMatch(str); m != nil {
		data = append(data, addData(m))
		for m != nil {
			m, _ = reg.FindNextMatch(m)
			if m != nil {
				data = append(data, addData(m))
			}
		}
	}
	return data
}

func ContentSha(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func TokensLength(str []string) int {
	length := 0
	for _, token := range str {
		length += TokenLength(token)
	}
	return length
}

func TokenLength(str string) int {
	reg := regexp.MustCompile(`[\w]+`)
	length := len(reg.FindAllString(str, -1))
	for i := 0; i < len(str); {
		_, size := utf8.DecodeRuneInString(str[i:])
		if size > 1 {
			length += 1
		}
		i += size
	}
	return length
}

func Header(headers map[string][]string, key string) string {
	if values, _ := headers[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func EbbedFile() string {
	embed, _ := web.BuildFS.ReadFile("build/js/embed.js")
	tmpl := `(function(){ {{.Ebbed}} }())`
	p := PromptTemplate(tmpl)
	//fmt.Println(string(embed))
	res, _ := p.Render(struct {
		Ebbed string
	}{
		Ebbed: string(embed),
	})
	return res
}

func WriteEnv(list map[string]string) {
	path := "./etc"
	env := path + "/.env"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return
		}
		godotenv.Write(list, env)
	} else {
		old, _ := godotenv.Read(env)
		for key, val := range list {
			old[key] = val
		}
		godotenv.Write(old, env)
	}
}

func Open(something string) {
	var open string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		open = "open"
	case "linux":
		open = "xdg-open"
	case "windows":
		open = "cmd"
		args = []string{"/c", "start"}
	default:
		fmt.Println(fmt.Sprintf("unknown OS, running on CPU: %s", runtime.GOOS))
		return
	}
	args = append(args, something)
	exec.Command(open, args...).Start()
}
