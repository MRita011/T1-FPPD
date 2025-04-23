# Jogo Concorrente em Go

## 🧭 Visão Geral

Este projeto expande um jogo de aventura existente, originalmente desenvolvido em Go com interface de texto, adicionando **elementos concorrentes e autônomos** ao mapa do jogo.

A movimentação básica do personagem e o carregamento do mapa já são fornecidos no código-base inicial, disponibilizado no repositório:  
```https://github.com/mvneves/fppd-jogo```

## 🎯 Objetivos

- Adicionar **três ou mais elementos** novos que se comportem de forma **concorrente**, utilizando goroutines.
- Implementar **interações reais** entre o personagem e os elementos do mapa.
- Desenvolver a função de **interação ativa do personagem** com o ambiente ao redor.
- Criar **mecânicas novas de jogo**, como obstáculos dinâmicos, armadilhas, inimigos ou NPCs interativos.
- Utilizar conceitos como **exclusão mútua**, **canais**, **select com múltiplos canais** e **timeouts** para comunicação entre threads.

## 🕹️ Como Jogar

- Use as teclas `W`, `A`, `S`, `D` para mover o personagem nas quatro direções.
- Use `E` para interagir com elementos próximos no mapa.
- Use `ESC` para encerrar o jogo.
- O personagem é representado por `☺`, e o mapa inicial é carregado a partir do arquivo `mapa.txt`.

## 🌍 Mapa

O mapa é uma matriz de 30x60 (modificável). Cada célula contém um caractere que representa um elemento, como:

- ` ` (espaço): célula vazia  
- `▤`: parede de tijolos  
- `♣`: vegetação  
- `☺`: personagem jogador  
- Outros símbolos: utilizados para os novos elementos implementados

## ⚙️ Novos Elementos Concorrentes

### 📦 Caixas Misteriosas (`■`)
Contêm tesouro, armadilha ou estão vazias. Possuem comportamento concorrente:
- Movimentam-se aleatoriamente a cada 20 segundos (timeout).
- Escutam canais para interação e decidem ação via `select`.
- Mudam de cor ao serem abertas, indicando seu conteúdo.

### 💰 Tesouros
Ocultos nas caixas misteriosas. Ao serem encontrados:
- Incrementam a contagem de vitórias.
- São removidos do mapa.
- Protegidos por `Mutex` durante modificação da matriz.

### 💣 Armadilhas
Também ocultas nas caixas misteriosas. Quando ativadas:
- Eliminam o jogador imediatamente.
- Disparam mensagens e encerram a partida.

### 🧙 NPC Guia (`🧙`)
- Inicia automaticamente em uma posição adjacente ao jogador.
- Segue o jogador dinamicamente.
- Fornece dicas em tempo real com base na distância:
  - 🔥 "Quente" → tesouro próximo
  - ❄️ "Frio" → armadilha próxima
  - 🌡️ "Morno" → objeto próximo
- Executa sua lógica em uma goroutine independente.

### 👾 Monstro (`¥`)
- Surge automaticamente após 30 segundos de jogo (`timeout`).
- Alterna entre movimentação aleatória e perseguição ao jogador.
- Rouba tesouros e pode encerrar a partida se coletar todos.
- Roda em uma goroutine dedicada, com mutex para coordenar movimentações.

> Todos os elementos acima são concorrentes, controlados por **goroutines**, e interagem com o mapa ou jogador via **canais**, **mutexes**, **selects** e **timeouts**.

## 🔄 Interação com o Personagem

A tecla `E` ativa a interação com elementos num raio próximo:
- Tesouro → coleta
- Armadilha → jogador é eliminado
- Prioridade é dada ao elemento mais próximo

Comunicação entre jogador e caixas ocorre via canal `chan bool`, garantindo **desacoplamento** e **segurança concorrente**.

## 🛠️ Compilação

### 🪟 Windows

Compilar com:

```cmd
go build -o jogo.exe
```

### ▶️ Como Executar

> Certifique-se de ter o arquivo `mapa.txt` com um mapa válido na raiz do projeto.
> Depois é só rodar no terminal:

- **Windows**:
  ```cmd
  /.jogo
  ```

## 🧑‍💻 Grupo

- Amanda Wilmsen: amanda.wilmsen@edu.pucrs.br  
- Killian D.B: killian.d@edu.pucrs.br  
- Luís Trein:  luis.trein@edu.pucrs.br  
- Maria Rita: m.ritarodrigues09@gmail.com  

## 📄 Relatório

O relatório em anexo descreve:

- Os elementos concorrentes planejados  
- Como interagem com o jogador  
- Comportamentos esperados  
- Abordagem de implementação usando goroutines, canais, mutexes e timeouts  

📄 [Link do Relatório no DOCS](https://docs.google.com/document/d/1BOLIXdguUHU_Q2kOid4UDEJaYej-GEtfyy_RirfOZLo/edit)