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
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			close(parar) // Fecha o canal para parar as goroutines quando o jogo terminar
			break
		}
	}
}
