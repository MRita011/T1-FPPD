// npc.go - Implementação do NPC Guian que segue o jogador
package main

import (
	"sync"
	"time"
)

// NPCGuian representa o estado do NPC guia
type NPCGuian struct {
	PosX, PosY int        // Posição atual do NPC
	Ativo      bool       // Indica se o NPC está ativo
	mu         sync.Mutex // Mutex para sincronização de acesso ao NPC
}

// Elemento visual do NPC Guian
var (
	NPC = Elemento{'🧙', CorVerde, CorPadrao, false} // Símbolo ♦ com cor verde
)

// Inicia o NPC em uma posição válida próxima ao jogador
func npcIniciar(jogo *Jogo) *NPCGuian {
	// Cria um novo NPC
	npc := &NPCGuian{
		Ativo: true,
	}

	// Encontra uma posição inicial válida para o NPC (próxima ao jogador)
	encontrarPosicaoInicial(jogo, npc)

	// Inicia a goroutine que controla o movimento do NPC
	go npcExecutar(jogo, npc)

	return npc
}

// Encontra uma posição válida para o NPC iniciar (próxima ao jogador)
func encontrarPosicaoInicial(jogo *Jogo, npc *NPCGuian) {
	// Direções possíveis para verificar (acima, direita, abaixo, esquerda)
	direcoes := [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	// Verifica cada direção
	for _, dir := range direcoes {
		dx, dy := dir[0], dir[1]
		nx, ny := jogo.PosX+dx, jogo.PosY+dy

		// Se a posição for válida, coloca o NPC lá
		if jogoPodeMoverPara(jogo, nx, ny) {
			npc.PosX = nx
			npc.PosY = ny
			return
		}
	}

	// Se nenhuma posição adjacente for válida, coloca o NPC em algum lugar por perto
	// Tenta em um raio maior
	for r := 2; r < 5; r++ {
		for x := -r; x <= r; x++ {
			for y := -r; y <= r; y++ {
				if x == 0 && y == 0 {
					continue // Pula a posição do jogador
				}

				nx, ny := jogo.PosX+x, jogo.PosY+y
				if jogoPodeMoverPara(jogo, nx, ny) {
					npc.PosX = nx
					npc.PosY = ny
					return
				}
			}
		}
	}
}

// Move o NPC em direção ao jogador
func npcMoverEmDirecaoAoJogador(jogo *Jogo, npc *NPCGuian) {
	npc.mu.Lock()
	defer npc.mu.Unlock()

	// Calcula a direção para o jogador
	dx := 0
	if jogo.PosX > npc.PosX {
		dx = 1
	} else if jogo.PosX < npc.PosX {
		dx = -1
	}

	dy := 0
	if jogo.PosY > npc.PosY {
		dy = 1
	} else if jogo.PosY < npc.PosY {
		dy = -1
	}

	// Tenta mover na direção horizontal ou vertical (prefere o movimento mais distante)
	nx, ny := npc.PosX, npc.PosY

	// Verifica diferença absoluta nas coordenadas
	diffX := abs(jogo.PosX - npc.PosX)
	diffY := abs(jogo.PosY - npc.PosY)

	// Tenta mover primeiro na direção com maior diferença
	if diffX > diffY {
		// Tenta mover horizontalmente primeiro
		if jogoPodeMoverPara(jogo, npc.PosX+dx, npc.PosY) {
			nx = npc.PosX + dx
		} else if dy != 0 && jogoPodeMoverPara(jogo, npc.PosX, npc.PosY+dy) {
			ny = npc.PosY + dy
		}
	} else {
		// Tenta mover verticalmente primeiro
		if jogoPodeMoverPara(jogo, npc.PosX, npc.PosY+dy) {
			ny = npc.PosY + dy
		} else if dx != 0 && jogoPodeMoverPara(jogo, npc.PosX+dx, npc.PosY) {
			nx = npc.PosX + dx
		}
	}

	// Atualiza a posição se houve movimento
	npc.PosX, npc.PosY = nx, ny
}

// Função helper para calcular o valor absoluto
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Goroutine que executa o comportamento do NPC
func npcExecutar(jogo *Jogo, npc *NPCGuian) {
	// Intervalo de atualização (X segundos)
	intervalo := 500 * time.Millisecond

	// Loop principal do NPC
	for npc.Ativo {
		// Move o NPC em direção ao jogador
		npcMoverEmDirecaoAoJogador(jogo, npc)

		// Pausa pelo intervalo definido
		time.Sleep(intervalo)
	}
}

// Desenha o NPC no mapa
func npcDesenhar(jogo *Jogo, npc *NPCGuian) {
	npc.mu.Lock()
	defer npc.mu.Unlock()

	// Sobrescreve temporariamente com o NPC
	interfaceDesenharElemento(npc.PosX, npc.PosY, NPC)
}
