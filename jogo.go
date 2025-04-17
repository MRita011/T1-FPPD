// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"math/rand"
	"sync"
	"os"
	"time"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem

}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      string       // mensagem para a barra de status
	Guian          *NPCGuian    // referência ao NPC guia
	FimDeJogo      bool         // indica se o jogador perdeu o jogo
	Tesouros       int
	Caixas         []*Caixa     // lista de caixas no mapa
	MutexMapa      *sync.Mutex // mutex para proteger o acesso ao mapa
}

// Elementos visuais do jogo
var (
	Personagem         = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Parede             = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao          = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio              = Elemento{' ', CorPadrao, CorPadrao, false}
	CaixaElemento      = Elemento{'■', CorAmarela, CorPadrao, false} // novo: caixa misteriosa
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo {
		UltimoVisitado: Vazio,
		MutexMapa: &sync.Mutex{},
	}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			case CaixaElemento.simbolo:
				e = CaixaElemento
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// utilizando o seed para gerar números aleatórios
	rand.Seed(time.Now().UnixNano())
	numCaixas := 10
	tipos := []TipoCaixa{VAZIA, TESOURO, ARMADILHA}

	// agora, vamos colocar as caixas em posições vazias aleatórias
	for colocadas := 0; colocadas < numCaixas; {
		x := rand.Intn(len(jogo.Mapa[0])) // escolhe uma coluna aleatória
		y := rand.Intn(len(jogo.Mapa))    //     "    "  linha aleatória

		if jogo.Mapa[y][x] == Vazio {
			jogo.Mapa[y][x] = CaixaElemento
			tipo := tipos[rand.Intn(len(tipos))] // escolhe um tipo aleatório
			caixa := &Caixa {
				X: x,
				Y: y,
				Tipo: tipo,
				Mapa: &jogo.Mapa,
				Mutex: jogo.MutexMapa,
				Interacao: make(chan bool),
			}
			
			caixa.Iniciar(jogo) // inicia a caixa
			jogo.Caixas = append(jogo.Caixas, caixa) // adiciona a caixa à lista de caixas
			colocadas++ // incrementa o contador de caixas colocadas
		}
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento
}

func interagir(jogo *Jogo) {
	jogo.MutexMapa.Lock()
	defer jogo.MutexMapa.Unlock()

	for _, caixa := range jogo.Caixas {
		if caixa.X == jogo.PosX && caixa.Y == jogo.PosY && !caixa.Removida {
			caixa.Interacao <- true
			break
		}
	}
}