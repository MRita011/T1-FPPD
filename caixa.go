package main

import (
	"math/rand"
	"sync"
	"time"
)

type TipoCaixa int 

const (
	VAZIA      TipoCaixa = 0
	TESOURO    TipoCaixa = 1
	ARMADILHA  TipoCaixa = 2
)

type Caixa struct {
	X, Y        int
	Tipo        TipoCaixa
	Mapa        *[][] Elemento 
	Mutex       *sync.Mutex
	Interacao   chan bool
	Interagindo bool
	Removida     bool
}

// iniciando uma goroutine para a caixa mudar de lugar
// a cada 20 segundos

func (c *Caixa) Iniciar(jogo *Jogo) {
	go func() {
		for !c.Removida {
			select {
				case <-time.After(20 * time.Second):
					c.mover();
				case <-c.Interacao:
					c.efeito(jogo)
					return
			}
		}
	}()
}

// movendo a caixa aleatoriamente
func (c *Caixa) mover() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	(*c.Mapa)[c.Y][c.X] = Vazio

	largura := len((*c.Mapa)[0])
	altura := len(*c.Mapa)

	for {
		novoX := rand.Intn(largura)
		novoY := rand.Intn(altura)

		if (*c.Mapa)[novoY][novoX] == Vazio {
			c.X = novoX
			c.Y = novoY
			(*c.Mapa)[novoY][novoX] = CaixaElemento
			break
		}
	}
}

// consequencias de cada tipo de caixa
func (c *Caixa) efeito(jogo *Jogo) {
    c.Mutex.Lock()
    defer c.Mutex.Unlock()

    // comportamento de cada caixa
    switch c.Tipo {
		case VAZIA:
			jogo.StatusMsg = "...CAIXA VAZIA!"
		
		case TESOURO:
			jogo.StatusMsg = "TESOURO ENCONTRADO!"
			jogo.Tesouros++
			exibirMensagemTesouros(jogo)
			
			// o jogador achou todos os tesouros?
			if jogo.Tesouros == 4 {
				jogo.StatusMsg = "Parabéns! Você encontrou todos os 4 tesouros!"
				jogo.FimDeJogo = true
			}
		
		case ARMADILHA:
			jogo.StatusMsg = "GAME OVER!"
			jogo.FimDeJogo = true
    }

    // caixa desaparecendo
    for i := 0; i < 10; i++ {  // a animação pode durar 10 ciclos
        time.Sleep(100 * time.Millisecond) // espera 100ms entre cada iteração

        if i%2 == 0 {
            (*c.Mapa)[c.Y][c.X] = Vazio
        } else {
            (*c.Mapa)[c.Y][c.X] = CaixaElemento
        }

        // Atualiza a tela
        interfaceDesenharJogo(jogo)
        interfaceAtualizarTela()
    }

    // após a animação, removemos a caixa permanentemente
    (*c.Mapa)[c.Y][c.X] = Vazio
    c.Removida = true
}
