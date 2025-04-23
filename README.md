# Jogo Concorrente em Go

## 🧭 Visão Geral

Este projeto expande um jogo de aventura existente, originalmente desenvolvido em Go com interface de texto, adicionando **elementos concorrentes e autônomos** ao mapa do jogo. A movimentação básica do personagem e o carregamento do mapa já são fornecidos no código-base inicial, disponibilizado no repositório:  
```https://github.com/mvneves/fppd-jogo```

## 🎯 Objetivos

- Adicionar **três ou mais elementos** novos que se comportem de forma **concorrente**, utilizando goroutines.
- Implementar **interações reais** entre o personagem e os elementos do mapa.
- Desenvolver a função de **interação ativa do personagem** com o ambiente ao redor.
- Criar **mecânicas novas de jogo**, como obstáculos dinâmicos, armadilhas, inimigos ou NPCs interativos.

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

### 💰 Tesouros (`💰`)
Itens colecionáveis ocultos no mapa. Ao interagir, são removidos e contam para a vitória. Comportamento gerenciado por thread.

### 💣 Armadilhas (`💣`)
Obstáculos letais. Se acionadas, eliminam o jogador. Podem ter comportamento dinâmico (ex: se mover ou alternar estado).

### 🧙 NPC Guia (`🧙`)
Personagem que segue o jogador e dá dicas em tempo real:
- 🔥 "Quente" → tesouro próximo
- ❄️ "Frio" → armadilha próxima
- 🌡️ "Morno" → algo próximo

Reza a lenda que todos os elementos rodam em threads independentes.

## 🔄 Interação com o Personagem

A tecla `E` ativa a interação com elementos num raio próximo:
- Tesouro → coleta
- Armadilha → jogador é eliminado
- Prioridade é dada ao elemento mais próximo

## 🧪 Testes e Exemplos

O arquivo `mapa.txt` contém cenários com os elementos para testes. O NPC e os objetos interagíveis respondem dinamicamente à posição do jogador.



## 🧑‍💻 Equipe

- Amanda Wilmsen: amanda.wilmsen@edu.pucrs.br
- Killian D.B: killian.d@edu.pucrs.br
- Luís Trein: email   
- Maria Rita: m.ritarodrigues09@gmail.com  

## 📁 Estrutura do Projeto

```
.
├── main.go
├── elementos.go
├── mapa.txt
├── utils.go
├── README.md
└── ...
```

## 📄 Relatório

O relatório anexo descreve:

- Os elementos concorrentes planejados
- Como interagem com o jogador
- Comportamentos esperados
- Abordagem de implementação usando goroutines e canais
  https://docs.google.com/document/d/1BOLIXdguUHU_Q2kOid4UDEJaYej-GEtfyy_RirfOZLo/edit?usp=drivesdk
