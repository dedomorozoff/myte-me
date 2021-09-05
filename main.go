// myte project main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//структура ответа рандомной шутки
type joke struct {
	Value string
}

func main() {

	//адреса API

	urlRandom := "https://api.chucknorris.io/jokes/random"
	urlCategory := "https://api.chucknorris.io/jokes/categories"
	urlFromCategory := "https://api.chucknorris.io/jokes/random?category="

	// Создаем флаги

	// для случайной шутки
	randomCmd := flag.NewFlagSet("random", flag.ExitOnError)

	//для шуток по категории
	dumpCmd := flag.NewFlagSet("dump", flag.ExitOnError)

	//парсим количество шуток для категории
	jokesNumber := dumpCmd.Int("n", 5, "n")

	//если нет аргументов, выходим
	if len(os.Args) < 2 {
		fmt.Println("usage: joker random [dump -n <value>]")
		os.Exit(1)
	}

	//если есть аргументы, выполняем нужные действия

	switch os.Args[1] {

	//рандомная шутка
	case "random":
		err := randomCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}

		randomJoke := getRandomJoke(urlRandom)

		//показываем рандомную шутку

		fmt.Println(randomJoke)

	//шутки из всех категорий

	case "dump":

		_ = dumpCmd.Parse(os.Args[2:])

		//получаем категории
		categories := getCategory(urlCategory)

		//проходимся по категориям
		for i := 0; i < len(categories); i++ {

			//берем массив из шуток для каждой категории
			jokesFromCategory := getJokesWithCategories(urlFromCategory, categories[i], *jokesNumber)

			//создаем файлы для категорий
			f, err := os.OpenFile(categories[i]+".txt", os.O_CREATE|os.O_RDWR, 0777)
			if err != nil {
				panic(err)
			}
			//очищаем файл
			err = f.Truncate(0)

			//сохраняем шутки в файлы

			for i := 0; i < len(jokesFromCategory); i++ {
				_, err = io.WriteString(f, fmt.Sprintln(jokesFromCategory[i]))
				if err != nil {
					panic(err)
				}
			}
			err = f.Close()
			fmt.Println("Jokes for category " + categories[i] + " saved.")
		}
		fmt.Println("All jokes are saved!")

		//ответ по-умолчанию в консоль
	default:

		fmt.Println("usage: joker random [dump -n <value>]")
		//выходим из программы
		os.Exit(1)
	}

}

//Функция возвращает случайную шутку в виде строки байт, на входе ссылка на шутку

func getFromUrl(url string) (result []byte) {

	//http соединение
	Client := http.Client{}
	//создаем новый запрос
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	//свой юзер-агент )
	req.Header.Set("User-Agent", "joker-bot")

	res, getErr := Client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(res.Body)
	}
	//получаем ответ с сервера
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	result = body

	//возвращаем
	return result
}

//Функция возвращает категории шуток в формате массива строк, на входе ссылка на категории

func getCategory(url string) []string {

	//берем категории

	body := getFromUrl(url)

	//декодим категории

	var arr []string
	_ = json.Unmarshal(body, &arr)

	return arr
}

//Функция возвращает случайную шутку в формате строки, на входе ссылка на шутку

func getRandomJoke(url string) string {

	//http соединение с таймаутом

	body := getFromUrl(url)
	//декодим шутку
	resultReady := joke{}
	jsonErr := json.Unmarshal(body, &resultReady)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	jokeResult := resultReady.Value

	//возвращаем шутку
	return jokeResult
}

//Функция возвращает случайные шутки с учетом категории в формате массива строк, на входе ссылка на шутку, название категории, количество шуток

func getJokesWithCategories(url string, category string, n int) (jokesWithCategories []string) {
	for i := 0; i < n; i++ {
		oneJoke := getRandomJoke(url + category)
		jokesWithCategories = append(jokesWithCategories, oneJoke)
	}
	return jokesWithCategories
}
