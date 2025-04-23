package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Monstro struct {
	X, Y             int
	Ativo            bool
	mu               sync.Mutex
	Velocidade       time.Duration
	TesourosRoubados int
}

func monstroNovo() *Monstro {
	return &Monstro{
		Velocidade: 2 * time.Second,
	}
}

func (m *Monstro) Iniciar(jogo *Jogo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Ativo = true
	encontrarPosicaoInicialMonstro(jogo, m)

	go m.comportamento(jogo)
}

// Atualizar é chamado pelo loop principal para avançar o estado do monstro
func (m *Monstro) Atualizar(jogo *Jogo) {
	if m.Ativo {
		m.mover(jogo)
		m.roubarTesouro(jogo)
	}
}

func encontrarPosicaoInicialMonstro(jogo *Jogo, m *Monstro) {
	// Posiciona longe do jogador
	for tentativas := 0; tentativas < 100; tentativas++ {
		x := rand.Intn(len(jogo.Mapa[0]))
		y := rand.Intn(len(jogo.Mapa))

		if jogo.Mapa[y][x] == Vazio &&
			calculaDistancia(jogo.PosX, jogo.PosY, x, y) > 10 {
			m.X, m.Y = x, y
			return
		}
	}
	// Se não encontrar posição ideal, coloca em qualquer lugar vazio
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			if elem == Vazio {
				m.X, m.Y = x, y
				return
			}
		}
	}
}

func (m *Monstro) comportamento(jogo *Jogo) {
	for m.Ativo && !jogo.FimDeJogo {
		m.mover(jogo)
		m.roubarTesouro(jogo)
		time.Sleep(m.Velocidade)
	}
}

func (m *Monstro) mover(jogo *Jogo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 50% de chance de mover na direção do jogador, 50% aleatório
	var dx, dy int
	if rand.Intn(2) == 0 {
		// Movimento inteligente em direção ao jogador
		if jogo.PosX > m.X {
			dx = 1
		} else if jogo.PosX < m.X {
			dx = -1
		}
		if jogo.PosY > m.Y {
			dy = 1
		} else if jogo.PosY < m.Y {
			dy = -1
		}
	} else {
		// Movimento aleatório
		dx = rand.Intn(3) - 1 // -1, 0 ou 1
		dy = rand.Intn(3) - 1
	}

	nx, ny := m.X+dx, m.Y+dy

	if ny >= 0 && ny < len(jogo.Mapa) &&
		nx >= 0 && nx < len(jogo.Mapa[0]) &&
		!jogo.Mapa[ny][nx].tangivel {

		m.X, m.Y = nx, ny
	}
}

func (m *Monstro) roubarTesouro(jogo *Jogo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Verifica se está perto do jogador (distância <= 1)
	if calculaDistancia(jogo.PosX, jogo.PosY, m.X, m.Y) <= 1 && jogo.Tesouros > 0 {
		jogo.Tesouros--
		m.TesourosRoubados++
		jogo.SetMessage(fmt.Sprintf("O monstro roubou um tesouro! (%d restantes)", jogo.Tesouros), 1*time.Minute) 

		if jogo.Tesouros <= 0 {
			jogo.SetMessage("GAME OVER!\nO monstro roubou TODOS os tesouros!", 30*time.Second)
			jogo.FimDeJogo = true
			return
		}
	}
}

func (m *Monstro) derrotar(jogo *Jogo) {
	if m.TesourosRoubados > 0 {
		jogo.Tesouros += m.TesourosRoubados
		jogo.SetMessage(fmt.Sprintf("Você recuperou %d tesouros!", m.TesourosRoubados), 3*time.Second)
		m.TesourosRoubados = 0
	}
	m.Ativo = false
}