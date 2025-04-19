// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"math/rand"
	"jogo/util"
	"sync"
	"os"
	"time"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune   // símbolo que vai aparecer no mapa
	cor      Cor    // cor do símbolo
	corFundo Cor    // cor do fundo
	tangivel bool   // se for true, não dá pra passar por cima
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      string       // mensagem para a barra de status
    MsgExpira      time.Time     // quando a mensagem expira
    MsgMutex       sync.Mutex    // protege o acesso às mensagens
	Guian          *NPCGuian    // referência ao NPC guia
	Monstro        *Monstro     // referência ao monstro
    MonstroAtivo   bool        // indica se o monstro está ativo 
	MonstroSpawn   time.Time    // quando o monstro vai aparecer
	FimDeJogo      bool         // indica se o jogador finalizou o jogo
	Tesouros       int          //quantidade de tesouros coletados
	Caixas         []*Caixa     // lista de caixas no mapa
	MutexMapa      *sync.Mutex  // mutex para proteger o acesso ao mapa
}

// Elementos visuais do jogo
var (
	Personagem           = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Monstro              = Elemento{'👹', CorVermelho, CorPadrao, true}
	Parede               = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao            = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio                = Elemento{' ', CorPadrao, CorPadrao, false}

	CaixaElemento        = Elemento{'■', CorAmarela, CorPadrao, true} // caixa fechada

	// para animar as caixas coloridinhas após abrir
	CaixaTesouroAberta   = Elemento{'■', CorVerde, CorPadrao, false}
	CaixaArmadilhaAberta = Elemento{'■', CorVermelho, CorPadrao, false}
	CaixaVaziaAberta     = Elemento{'■', CorCinzaEscuro, CorPadrao, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{
		UltimoVisitado: Vazio,
		MutexMapa:      &sync.Mutex{},
        MonstroSpawn:   time.Now().Add(30 * time.Second), // monstro aparece após 30 segundos
	}
}

func atualizarJogo(jogo *Jogo) {
    // Spawn do monstro
    if !jogo.MonstroAtivo && time.Now().After(jogo.MonstroSpawn) {
        jogo.Monstro = monstroNovo()
        jogo.Monstro.Iniciar(jogo)
        jogo.MonstroAtivo = true
        jogo.SetMessage("Um monstro apareceu no mapa!\nCorra para pegar todos os tesouros antes que ele roube...", 5*time.Second)
    }
    
    // Atualizar monstro se estiver ativo
    if jogo.MonstroAtivo {
        jogo.Monstro.Atualizar(jogo)
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

	// coloca o seed pra gerar números aleatórios diferentes toda vez que o jogo é iniciado
	rand.Seed(time.Now().UnixNano())
	numCaixas := 10 // número de caixas pra espalhar no mapa
	tipos := []TipoCaixa{VAZIA, TESOURO, ARMADILHA}

	// agora espalha as caixas em lugares aleatórios que estão vazios
	for colocadas := 0; colocadas < numCaixas; {
		x := rand.Intn(len(jogo.Mapa[0])) // pega coluna aleatória
		y := rand.Intn(len(jogo.Mapa))    // pega linha aleatória

		if jogo.Mapa[y][x] == Vazio {
			jogo.Mapa[y][x] = CaixaElemento
			tipo := tipos[rand.Intn(len(tipos))] // escolhe um tipo de caixa (aleatoriamente)
			caixa := &Caixa{
				X:          x,
				Y:          y,
				Tipo:       tipo,
				Mapa:       &jogo.Mapa,
				Mutex:      jogo.MutexMapa,
				Interacao:  make(chan bool),
			}

			caixa.Iniciar(jogo) // inicia a caixa
			jogo.Caixas = append(jogo.Caixas, caixa) // adiciona na lista de caixas
			colocadas++ // marca que colocou uma
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

	// Verifica se o elemento de destino é uma caixa (não pode passar por cima)
	if jogo.MonstroAtivo && jogo.Monstro.X == x && jogo.Monstro.Y == y {
        return false
    }

	// liberado pra andar
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

// versão especial pra impedir andar por cima das caixas
func jogoMoverComCaixa(jogo *Jogo, x, y int) bool {
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// checa se tem alguma caixa bloqueando a passagem
	for _, caixa := range jogo.Caixas {
		if caixa.X == x && caixa.Y == y && !caixa.Removida {
			jogo.StatusMsg = "Uma caixa bloqueia o caminho!"
			return false
		}
	}

	if jogo.Mapa[y][x].tangivel {
		return false
	}
	return true
}

// permite o jogador interagir com caixas que estão até 1 célula de distância
func interagir(jogo *Jogo) {
	jogo.MutexMapa.Lock() // trava o mapa pra ninguém mexer enquanto interage
	defer jogo.MutexMapa.Unlock()

	for _, caixa := range jogo.Caixas {
		// checa se a caixa está próxima e não foi removida
		if !caixa.Removida && util.Abs(caixa.X-jogo.PosX) <= 1 && util.Abs(caixa.Y-jogo.PosY) <= 1 {
			caixa.Interacao <- true // manda o sinal pra caixa abrir
			jogo.StatusMsg = "Você interagiu com a caixa!"
			caixa.Removida = true // marca que a caixa foi removida
			break
		}
	}
}

func (j *Jogo) SetMessage(msg string, duration time.Duration) {
    j.MsgMutex.Lock()
    defer j.MsgMutex.Unlock()
    
    j.StatusMsg = msg
    j.MsgExpira = time.Now().Add(duration)
}

func (j *Jogo) GetMessage() string {
    j.MsgMutex.Lock()
    defer j.MsgMutex.Unlock()
    
    if time.Now().After(j.MsgExpira) {
        return ""
    }
    return j.StatusMsg
}