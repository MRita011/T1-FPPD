package main

import (
	"os"
	"time"
)

// Função para iniciar a renderização periódica do jogo
func iniciarRenderizador(jogo *Jogo, parar <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond) // redesenha a cada 100ms
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				interfaceDesenharJogo(jogo)
			case <-parar:
				return
			}
		}
	}()
}

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Cria canal para parar as goroutines (NPC e renderizador)
	parar := make(chan struct{})

	// Inicializa o NPC
	jogo.Guian = npcIniciar(&jogo)

	// Inicia a renderização contínua do jogo
	iniciarRenderizador(&jogo, parar)

	// Loop principal de entrada
	for !jogo.FimDeJogo {
		// Atualiza o estado do jogo (monstro, npc, etc)
		atualizarJogo(&jogo)

		// Processa entrada do usuário
		evento := interfaceLerEventoTeclado()
		if evento.Tipo == "sair" {
			break
		}

		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}

		// Pequena pausa para evitar uso excessivo da CPU
		time.Sleep(50 * time.Millisecond)
	}

	// Fecha o canal para parar as goroutines
	close(parar)

	// Mostra mensagem final por mais tempo se o jogo terminou
	if jogo.FimDeJogo {
		interfaceDesenharJogo(&jogo)
		time.Sleep(5 * time.Second)
	}
}