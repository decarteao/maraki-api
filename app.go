package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maraki/scrapper"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// novo roteador
var routes = mux.NewRouter()

// manipular de acesso
func ValidarRequisicao(c *http.Request) bool {
	return c.Header.Get("maraki") == "online.helio3.marakitv"
}

// handlers das rotas
func HandlerIndex(w http.ResponseWriter, c *http.Request) {
	// para a pagina inicial do site
	w.Header().Set("Content-Type", "application/json")

	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		json.NewEncoder(w).Encode(scrapper.Index())
	}
}
func HandlerFilmes(w http.ResponseWriter, c *http.Request) {
	// para a pagina filmes
	w.Header().Set("Content-Type", "application/json")

	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		page := mux.Vars(c)["page"]
		page_int, err := strconv.Atoi(page)

		if err != nil {
			page = "1"
		} else if page_int < 1 {
			page = "1"
		}

		json.NewEncoder(w).Encode(scrapper.ListarFilmes(page))
	}
}
func HandlerSeries(w http.ResponseWriter, c *http.Request) {
	// para a pagina de series
	w.Header().Set("Content-Type", "application/json")

	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		page := mux.Vars(c)["page"]
		page_int, err := strconv.Atoi(page)

		if err != nil {
			page = "1"
		} else if page_int < 1 {
			page = "1"
		}

		json.NewEncoder(w).Encode(scrapper.ListarSeries(page))
	}
}
func HandlerListarCategoria(w http.ResponseWriter, c *http.Request) {
	// para algum categoria escolhido
	w.Header().Set("Content-Type", "application/json")
	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		json.NewEncoder(w).Encode(scrapper.ListarCategoria(mux.Vars(c)["categoria"]))
	}
}
func HandlerPesquisa(w http.ResponseWriter, c *http.Request) {
	// para coisas pesquisadas
	w.Header().Set("Content-Type", "application/json")
	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		json.NewEncoder(w).Encode(scrapper.GetResults(c.URL.Query().Get("q")))
	}
}
func HandlerListarTemporada(w http.ResponseWriter, c *http.Request) {
	// listar episodios de uma temporada
	w.Header().Set("Content-Type", "application/json")

	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		// retornar
		vars := mux.Vars(c)
		json.NewEncoder(w).Encode(scrapper.ListarTemporada(vars["id"], vars["temporada"]))
	}
}
func HandlerListarCategorias(w http.ResponseWriter, c *http.Request) {
	// listar categorias
	w.Header().Set("Content-Type", "application/json")
	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		// retornar
		json.NewEncoder(w).Encode(scrapper.Generos())
	}
}
func HandlerAbrirFilme(w http.ResponseWriter, c *http.Request) {
	// retorna os links dos videos
	w.Header().Set("Content-Type", "application/json")

	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		// retornar
		json.NewEncoder(w).Encode(scrapper.AbrirFilme(mux.Vars(c)["id"]))
	}
}
func HandlerAbrirSerie(w http.ResponseWriter, c *http.Request) {
	// retorna os links dos videos
	w.Header().Set("Content-Type", "application/json")
	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		// retornar
		json.NewEncoder(w).Encode(scrapper.AbrirSerie(mux.Vars(c)["id"]))
	}
}

func HandlerAssistirFilme(w http.ResponseWriter, c *http.Request) {
	// retorna os links dos videos
	w.Header().Set("Content-Type", "application/json")
	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		// get url
		url := c.URL.Query().Get("url")

		// retornar
		json.NewEncoder(w).Encode(scrapper.GetFilmeLinkVideo(url))
	}
}
func HandlerAssistirEpisodio(w http.ResponseWriter, c *http.Request) {
	// retorna os links de video do episodio
	w.Header().Set("Content-Type", "application/json")
	// validar
	if !ValidarRequisicao(c) {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid resquest"})
	} else {
		// retornar
		json.NewEncoder(w).Encode(scrapper.GetSerieLinkVideo(mux.Vars(c)["episodio"]))
	}
}

func main() {
	// rotas
	routes.HandleFunc("/", HandlerIndex)

	// listar filmes
	routes.HandleFunc("/filmes", HandlerFilmes)
	// listar filmes pela indexacao
	routes.HandleFunc("/filmes/{page}", HandlerFilmes)

	// listar series
	routes.HandleFunc("/series", HandlerSeries)
	// listar series com indexacao
	routes.HandleFunc("/series/{page}", HandlerSeries)

	// pesquisar algo no parametro: ?q=aqui+fica+a+busca
	routes.HandleFunc("/busca", HandlerPesquisa)

	// listar categorias de filmes
	routes.HandleFunc("/listar/categorias", HandlerListarCategorias)
	// quando escolhido uma categoria
	routes.HandleFunc("/categoria/{categoria}", HandlerListarCategoria)

	// visualizar a pagina do filme ou serie
	routes.HandleFunc("/serie/{id}", HandlerAbrirSerie)
	routes.HandleFunc("/filme/{id}", HandlerAbrirFilme)

	// listar episodios
	routes.HandleFunc("/serie/{id}/{temporada}", HandlerListarTemporada)

	// ver episodio
	routes.HandleFunc("/serie/{id}/{temporada}/{episodio}", HandlerAssistirEpisodio)

	// ver filme ao mandar o link do player em: ?url=urlencode link
	routes.HandleFunc("/filme/{id}/ver", HandlerAssistirFilme)

	// fmt.Println("Webservice...")
	// http.ListenAndServe(":80", routes)

	port := "8080" // Porta padrão para Railway
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, routes))
}
