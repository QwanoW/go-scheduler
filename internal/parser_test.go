package internal

import (
	"regexp"
	"strings"
	"testing"
)

// wildCardToRegexp преобразует строку с подстановочными знаками в строку с регулярным выражением.
func wildCardToRegexp(pattern string) string {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {
		// Заменяем * на .*
		if i > 0 {
			result.WriteString(".*")
		}
		// Экранируем любые специальные символы в литерале.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return result.String()
}

func Test_findLinks(t *testing.T) {
	type args struct {
		url         string
		parseTarget string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			"test1",
			args{"https://sustec.ru/raspisanie-ptk-ul-gagarina-7/", "mtli_doc"},
			map[string]string{
				"*1,2 курсы": viewerTemplate + "https://sustec.ru*",
				"*3,4 курсы": viewerTemplate + "https://sustec.ru*",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindLinks(tt.args.url, tt.args.parseTarget)
			if (err != nil) != tt.wantErr {
				t.Errorf("findLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("findLinks() = %v, want %v", got, tt.want)
				return
			}

			// Сравниваем каждую пару ключ-значение с помощью регулярных выражений
			for k, v := range tt.want {
				// Преобразуем подстановочные знаки в регулярные выражения
				kRegexp := wildCardToRegexp(k)
				vRegexp := wildCardToRegexp(v)
				// Проверяем, есть ли такой ключ в полученном значении
				kMatch, _ := regexp.MatchString(kRegexp, k)
				vMatch, _ := regexp.MatchString(vRegexp, v)
				if !kMatch {
					t.Errorf("Missing key %s in findLinks() = %v, want %v", k, got, tt.want)
					return
				}
				if !vMatch {
					t.Errorf("Missing value %s in findLinks() = %v, want %v", v, got, tt.want)
					return
				}
			}
		})
	}
}
