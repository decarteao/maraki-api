package scrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// variaveis globais
const URL = "https://assistir.biz"

var client = &http.Client{
	Timeout: 10 * time.Second,
}

// structs do scrapper
type Card struct {
	Title    string `json:"title"`
	Image    string `json:"image"`
	Category string `json:"category"`
	Ano      string `json:"ano"`
	Link     string `json:"link"`
	Rank     string `json:"rank"`
}
type CardFullFilme struct {
	Title          string `json:"title"`
	Image          string `json:"image"`
	Category       string `json:"category"`
	Ano            string `json:"ano"`
	Link1          string `json:"link_1"`
	Link2          string `json:"link_2"`
	Link3          string `json:"link_3"`
	Rank           string `json:"rank"`
	Duration       string `json:"duration"`
	Classification string `json:"classification"`
	Description    string `json:"description"`
}
type CardFullSerie struct {
	Title       string            `json:"title"`
	Image       string            `json:"image"`
	Category    string            `json:"category"`
	Ano         string            `json:"ano"`
	Duration    string            `json:"duration"`
	Description string            `json:"description"`
	Temporadas  map[string]string `json:"temporadas"`
}
type IndexPage struct {
	Destaques      []Card `json:"destaques"`
	MaisAssistidos []Card `json:"mais_assistidos"`
	UltAdicionados []Card `json:"ult_adicionados"`
	Series         []Card `json:"series"`
}
type JSONAssistir struct {
	// {"id":"2837","hd":"1","token":"75be1128463c6e0a4395e6d23b60aadc","serie_ep":"1-1.mp4","dir_path":"acasadodragao","hls":"0"}
	Id      string `json:"id"`
	Hd      string `json:"hd"`
	Token   string `json:"token"`
	SerieEp string `json:"serie_ep"`
	DirPath string `json:"dir_path"`
	HLS     string `json:"hls"`
}

// funcao de setar os headers
func SetHeadersClient(req *http.Request, headers *map[string]string) {
	// headers padroes
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")

	// headers personalizaveis
	for key, val := range *headers {
		req.Header.Set(key, val)
	}
}

// funcoes principais do scrapper
func FormatarLink(url string) string {
	if url[:3] == "//a" {
		url = fmt.Sprintf("https:%s", url)
	}
	return url
}

func GetCard(element *goquery.Selection, tipo rune) Card {
	var card Card

	link, _ := element.Find("a[onclick]").Attr("href")

	img, exists := element.Find("img").Attr("src")
	if !exists || strings.Contains(img, "poster_default") || strings.Contains(img, "_filter(blur)") {
		img = element.Find("img").AttrOr("data-src", "https://assistir.biz/assets/img/poster_default.jpg")
	}

	if tipo == 'i' {
		// index
		card = Card{
			Title:    strings.Trim(element.Find(".card__title").Text(), " \n\r\t"),
			Image:    strings.Trim(img, " \n\r\t"),
			Category: strings.Trim(element.Find(".card__category > .span-category").Text(), " \n\r\t"),
			Ano:      strings.Trim(element.Find(".card__category > .span-year").Text(), " \n\r\t"),
			Link:     strings.Trim(link, " \n\r\t"),
			Rank:     strings.Trim(element.Find(".card__rate").Text(), " \n\r\t"),
		}
	} else if tipo == 'b' {
		// busca
		cat_ano := strings.Split(element.Find(".card__category").Text(), ",")
		var (
			cat string
			ano string
		)

		if len(cat_ano) == 2 {
			cat = cat_ano[0]
			ano = cat_ano[1]
		} else {
			cat = ""
			ano = cat_ano[0]
		}

		card = Card{
			Title:    strings.Trim(element.Find(".card__title").Text(), " \n\r\t"),
			Image:    strings.Trim(img, " \n\r\t"),
			Category: strings.Trim(cat, " \n\r\t"),
			Ano:      strings.Trim(ano, " \n\r\t"),
			Link:     strings.Trim(link, " \n\r\t"),
			Rank:     strings.Trim(element.Find(".card__rate").Text(), " \n\r\t"),
		}
	} else if tipo == 'f' {
		// filmes
		cat_ano := strings.Split(element.Find(".card__category").Text(), ",")
		var (
			cat string
			ano string
		)

		if len(cat_ano) == 2 {
			cat = cat_ano[0]
			ano = cat_ano[1]
		} else {
			cat = ""
			ano = cat_ano[0]
		}

		card = Card{
			Title:    strings.Trim(element.Find(".card__title").Text(), " \n\r\t"),
			Image:    strings.Trim(img, " \n\r\t"),
			Category: strings.Trim(cat, " \n\r\t"),
			Ano:      strings.Trim(ano, " \n\r\t"),
			Link:     strings.Trim(link, " \n\r\t"),
			Rank:     strings.Trim(element.Find(".card__rate").Text(), " \n\r\t"),
		}
	} else if tipo == 's' {
		// series
		cat_ano := strings.Split(strings.Trim(strings.Trim(element.Find(".card__category").Text(), " "), "\n"), "\n")
		var (
			cat string
			ano string
		)

		if len(cat_ano) == 2 {
			cat = cat_ano[0]
			ano = cat_ano[1]
		} else {
			cat = ""
			ano = cat_ano[0]
		}

		card = Card{
			Title:    strings.Trim(element.Find(".card__title").Text(), " \n\r\t"),
			Image:    strings.Trim(img, " \n\r\t"),
			Category: strings.Trim(cat, " \n\r\t"),
			Ano:      strings.Trim(ano, " \n\r\t"),
			Link:     strings.Trim(link, " \n\r\t"),
			Rank:     strings.Trim(element.Find(".card__rate").Text(), " \n\r\t"),
		}
	} else if tipo == 'l' {
		// listagem de series ou filmes
		cat_ano := strings.Split(element.Find(".card__category").Text(), ",")
		var (
			cat string
			ano string
		)

		if len(cat_ano) == 2 {
			cat = cat_ano[0]
			ano = cat_ano[1]
		} else {
			cat = ""
			ano = cat_ano[0]
		}

		card = Card{
			Title:    strings.Trim(element.Find(".card__title").Text(), " \n\r\t"),
			Image:    strings.Trim(img, " \n\r\t"),
			Category: strings.Trim(cat, " \n\r\t"),
			Ano:      strings.Trim(ano, " \n\r\t"),
			Link:     strings.Trim(link, " \n\r\t"),
			Rank:     strings.Trim(element.Find(".card__rate").Text(), " \n\r\t"),
		}
	} else if tipo == 'g' {
		// gender - categoria
		cat := element.Find(".card__category > .span-category").Text()
		ano := element.Find(".card__category > .span-year").Text()

		card = Card{
			Title:    strings.Trim(element.Find(".card__title").Text(), " \n\r\t"),
			Image:    strings.Trim(img, " \n\r\t"),
			Category: strings.Trim(cat, " \n\r\t"),
			Ano:      strings.Trim(ano, " \n\r\t"),
			Link:     strings.Trim(link, " \n\r\t"),
			Rank:     strings.Trim(element.Find(".card__rate").Text(), " \n\r\t"),
		}
	}

	return card
}

func Index() interface{} {
	// fazer requisicao
	req, _ := http.NewRequest("GET", URL, nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	// iniciar o scraping
	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	// data geral
	data_geral := IndexPage{}

	// pegar os destaques
	soup.Find("h1[class*=\"home__title\"] > b").Parent().Parent().Parent().Find("div.col-12 > div[class*=\"owl-carousel\"] > div[id][class*=\"card card\"]").Each(func(i int, s *goquery.Selection) {
		data_geral.Destaques = append(data_geral.Destaques, GetCard(s, 'i'))
	})

	// pegar os mais assistidos
	soup.Find("button[class*=\"section__nav\"][type=\"button\"][data-nav=\"#owl_assistidos\"]").Parent().Parent().Parent().Parent().Find("div.col-12 > div[class*=\"owl-carousel\"] > div[id][class*=\"card\"]").Each(func(i int, s *goquery.Selection) {
		data_geral.MaisAssistidos = append(data_geral.MaisAssistidos, GetCard(s, 'i'))
	})

	// ultimos assistidos
	soup.Find("button[class*=\"section__nav\"][type=\"button\"][data-nav=\"#owl_ultimos\"]").Parent().Parent().Parent().Parent().Find("div.col-12 > div[class*=\"owl-carousel\"] > div[id][class*=\"card\"]").Each(func(i int, s *goquery.Selection) {
		data_geral.UltAdicionados = append(data_geral.UltAdicionados, GetCard(s, 'i'))
	})

	// series
	soup.Find("button[class*=\"section__nav\"][type=\"button\"][data-nav=\"#owl_series\"]").Parent().Parent().Parent().Parent().Find("div.col-12 > div[class*=\"owl-carousel\"] > div[id][class*=\"card\"]").Each(func(i int, s *goquery.Selection) {
		data_geral.Series = append(data_geral.Series, GetCard(s, 'i'))
	})

	return data_geral
}

func Generos() interface{} {
	// fazer requisicao
	req, _ := http.NewRequest("GET", URL, nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	generos := map[string]string{}
	soup.Find("ul[aria-labelledby=\"dropdownMenuCatalog\"] a[href*=\"/categoria/\"]").Each(func(i int, s *goquery.Selection) {
		generos[strings.Trim(s.Text(), " \n\r")] = strings.Trim(s.AttrOr("href", "#"), " \n\r")
	})

	return generos
}

func GetResults(search string) interface{} {
	// para esquema de busca
	search = strings.ReplaceAll(strings.Trim(search, " "), " ", "+")
	url_schema := fmt.Sprintf("%s/busca?q=%s", URL, search)

	// fazer requisicao
	req, _ := http.NewRequest("GET", url_schema, nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	data_geral := []Card{}

	soup.Find("div[class='card']").Each(func(idx int, el *goquery.Selection) {
		data_geral = append(data_geral, GetCard(el, 'b'))
	})

	return data_geral
}

func ListarCategoria(categoria string) interface{} {
	// quand escolhe um genero especifico

	// fazer requisicao
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/categoria/%s", URL, categoria), nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	data_geral := []Card{}

	soup.Find("div[class='card']").Each(func(idx int, el *goquery.Selection) {
		data_geral = append(data_geral, GetCard(el, 'g'))
	})

	return data_geral
}

func ListarFilmes(page string) interface{} {
	// fazer requisicao
	data := url.Values{}
	data.Set("page", page)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/listaFilmes", URL), strings.NewReader(data.Encode()))
	SetHeadersClient(req, &map[string]string{"x-requested-with": "XMLHttpRequest", "content-type": "application/x-www-form-urlencoded; charset=UTF-8"})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	// d, _ := io.ReadAll(resp.Body)
	// os.WriteFile("test.html", d, 0777)

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	data_geral := map[string]interface{}{"page": "", "next": "", "prev": "", "cards": []Card{}}

	soup.Find("div[class='card']").Each(func(idx int, el *goquery.Selection) {
		data_geral["cards"] = append(data_geral["cards"].([]Card), GetCard(el, 'f'))
	})

	paginacao := soup.Find("ul[id=\"paginacao\"][class=\"paginator\"]")
	current_page, _ := strconv.Atoi(paginacao.Find("li[class=\"paginator__item paginator__item--active\"]").Text())

	prev_page, next_page := "", ""

	if len(paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page-1)).Text()) > 0 {
		prev_page = paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page-1)).Text()
	}
	if len(paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page+1)).Text()) > 0 {
		next_page = paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page+1)).Text()
	}

	data_geral["page"] = strconv.Itoa(current_page)
	data_geral["prev"] = prev_page
	data_geral["next"] = next_page

	return data_geral
}

func ListarSeries(page string) interface{} {
	// fazer requisicao
	data := url.Values{}
	data.Set("page", page)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/series", URL), strings.NewReader(data.Encode()))
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	data_geral := map[string]interface{}{"page": "", "next": "", "prev": "", "cards": []Card{}}

	soup.Find("div[class='card']").Each(func(idx int, el *goquery.Selection) {
		data_geral["cards"] = append(data_geral["cards"].([]Card), GetCard(el, 's'))
	})

	paginacao := soup.Find("ul[id=\"paginacao\"][class=\"paginator\"]")
	current_page, _ := strconv.Atoi(paginacao.Find("li[class=\"paginator__item paginator__item--active\"]").Text())
	if current_page == 0 {
		current_page = 1
	}

	prev_page, next_page := "", ""

	if len(paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page-1)).Text()) > 0 {
		prev_page = paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page-1)).Text()
	}
	if len(paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page+1)).Text()) > 0 {
		next_page = paginacao.Find(fmt.Sprintf("a[data-page=\"%v\"]", current_page+1)).Text()
	}

	data_geral["page"] = strconv.Itoa(current_page)
	data_geral["prev"] = prev_page
	data_geral["next"] = next_page

	return data_geral
}

func GetFilmeLinkVideo(link_player string) interface{} {
	req, _ := http.NewRequest("GET", link_player, nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	videos := map[int]string{}

	soup.Find("video source[src][size]").Each(func(i int, s *goquery.Selection) {
		key, _ := strconv.Atoi(s.AttrOr("size", "0"))
		videos[key] = FormatarLink(s.AttrOr("src", "#"))
	})

	return videos
}

func AbrirFilme(id string) interface{} {
	// retorna: os dados do CardFull e sugestoes de outros filmes | e trailer

	// fazer requisicao
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/filme/%s", URL, id), nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	data_geral := map[string]interface{}{"card": CardFullFilme{}, "trailer": nil, "sugestoes": []Card{}}

	// pegar os links de visualizacao
	players := []string{}

	soup.Find("div[id*=\"player-\"]").Each(func(i int, s *goquery.Selection) {
		players = append(players, FormatarLink(s.Find("iframe").AttrOr("src", "#")))
	})

	var link_1, link_2, link_3 string

	if len(players) > 0 {
		link_1 = players[0]
		if len(players) > 1 {
			link_2 = players[1]
			if len(players) > 2 {
				link_3 = players[2]
			}
		}
	} else {
		return map[string]string{"error": "not found"}
	}

	// pegar o card
	title := strings.Trim(soup.Find("section[class=\"section section--details\"] h1.section__title").Text(), " ")
	rank := strings.Trim(soup.Find("section[class=\"section section--details\"] .card__rate").Text(), " ")

	details := soup.Find("section[class=\"section section--details\"] .container > .row .card--details > .row")

	ano := strings.Split(strings.Trim(details.Find(".card__meta > li i[class*=\"fa-calendar-days\"]").Parent().Parent().Text(), " \n\r\t"), "\n")[1]
	ano = strings.Trim(ano, " \t\n")

	dur := strings.Split(strings.Trim(details.Find(".card__meta > li i[class*=\"fa-timer\"]").Parent().Text(), " \n\r\t"), ":")[1]
	dur = strings.Trim(dur, " ")

	indicacao := details.Find(".card__meta > li #class_indicativa").Text()
	description := strings.Trim(details.Find(".card__description").Text(), " \n\t\r")

	data_geral["card"] = CardFullFilme{
		Title:          title,
		Image:          details.Find("img[data-src]").AttrOr("data-src", "https://assistir.biz/assets/img/poster_default.webp"),
		Category:       strings.Trim(details.Find(".card__meta > li a[href]:not([style])").Text(), " "),
		Ano:            ano,
		Rank:           rank,
		Duration:       dur,
		Classification: indicacao,
		Description:    description,
		Link1:          link_1,
		Link2:          link_2,
		Link3:          link_3,
	}

	// sugestoes
	soup.Find("div[class='card']").Each(func(idx int, el *goquery.Selection) {
		data_geral["sugestoes"] = append(data_geral["sugestoes"].([]Card), GetCard(el, 'l'))
	})

	// trailer
	data_geral["trailer"] = soup.Find("lite-youtube[videoid]").AttrOr("videoid", "#")

	return data_geral
}

func GetSerieLinkVideo(idvideo string) interface{} {
	// aqui busca o link do video pelo id
	p := url.Values{}
	p.Set("id", idvideo)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/getepisodio", URL), strings.NewReader(p.Encode()))

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	headers := map[string]string{
		"origin": URL, "x-requested-with": "XMLHttpRequest", "content-type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	SetHeadersClient(req, &headers)

	// send post
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	buff, _ := io.ReadAll(resp.Body)

	if strings.Contains(string(buff), "tá fazendo o quê aqui?") {
		fmt.Println("error: nao achou o link")
		return map[string]string{"error": "nao achou o link"}
	}

	var data_json JSONAssistir

	err = json.NewDecoder(bytes.NewBuffer(buff)).Decode(&data_json)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	buff = nil

	link := fmt.Sprintf("https://assistir.biz/playserie/%v/%s", idvideo, data_json.Token)

	return map[string]string{"video": link}
}

func ListarTemporada(id, temporada string) interface{} {
	// retorna: os episodios de uma temporada

	// fazer requisicao
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/serie/%s/%s", URL, id, temporada), nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	ids := map[string]string{}

	soup.Find("tr[onclick*=\"reloadVideoSerie(\"]").Each(func(i int, s *goquery.Selection) {
		id := s.AttrOr("onclick", "reloadVideoSerie(2837, 'f576078a6d201711ba4834c58f138b20')")
		id = strings.Trim(strings.Split(strings.Split(id, ",")[0], "(")[1], " ")

		pos := s.Find("th").First().Text()
		title := strings.Trim(s.Find("th").Last().Text(), " \t\n\r")

		ids[id] = fmt.Sprintf("%s - %s", pos, title)
	})

	return ids
}

func AbrirSerie(id string) interface{} {
	// retorna: os dados do CardFull e sugestoes | e trailer
	// fazer requisicao
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/serie/%s", URL, id), nil)
	SetHeadersClient(req, &map[string]string{})

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	defer resp.Body.Close()

	soup, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
		return map[string]string{"error": err.Error()}
	}

	data_geral := map[string]interface{}{"card": CardFullSerie{}, "trailer": nil, "sugestoes": []Card{}}

	// pegar o card
	title := strings.Trim(soup.Find("section[class=\"section section--details\"] h1.section__title").Text(), " ")

	var categoria, ano, dur, descricao string

	details := soup.Find("section[class=\"section section--details\"] .container > .row .card--details > .row")
	descricao = strings.Trim(details.Find(".card__description").Text(), " \n\t\r")

	details.Find(".card__meta > li").Each(func(i int, s *goquery.Selection) {
		key_val := strings.Split(s.Text(), ":")
		key, val := key_val[0], key_val[1]
		val = strings.Trim(val, " \n\t\r,")

		if strings.Contains(key, "Gênero") {
			categoria = val
		}
		if strings.Contains(key, "Ano de lançamento") {
			ano = val
		}
		if strings.Contains(key, "Duração") {
			dur = val
		}
	})

	// pegar as temporadas
	temps := map[string]string{}

	soup.Find("div[class=\"card\"][id*=\"temporada-\"]").Each(func(i int, s *goquery.Selection) {
		temp_name := strings.Split(s.AttrOr("id", "1"), "-")[1]
		link := s.Find("a.card__play").AttrOr("href", "#")
		temps[temp_name] = link
	})

	data_geral["card"] = CardFullSerie{
		Title:       title,
		Image:       details.Find("img[data-src]").AttrOr("data-src", "https://assistir.biz/assets/img/poster_default.webp"),
		Category:    categoria,
		Ano:         ano,
		Duration:    dur,
		Description: descricao,
		Temporadas:  temps,
	}

	// sugestoes
	soup.Find("div[class='card']").Each(func(idx int, el *goquery.Selection) {
		if !strings.Contains(el.AttrOr("id", "temporada-"), "temporada-") {
			data_geral["sugestoes"] = append(data_geral["sugestoes"].([]Card), GetCard(el, 'l'))
		}
	})

	// trailer
	data_geral["trailer"] = soup.Find("lite-youtube[videoid]").AttrOr("videoid", "#")

	return data_geral
}
