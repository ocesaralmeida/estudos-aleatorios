# 1. Processos: O Ambiente de Execução

## Definição:
- Um processo é uma instância de um programa em execução. O sistema operacional cria um processo quando um programa é iniciado, fornecendo os recursos necessários para sua execução, como espaço de endereçamento de memória, descritores de arquivos, variáveis de ambiente, etc.

## Características:

- Cada processo tem seu próprio espaço de memória, isolado de outros processos.

- A comunicação entre processos (IPC) é mais complexa e custosa, pois não compartilham memória.

- O sistema operacional gerencia os processos, atribuindo tempo de CPU e outros recursos.

## Exemplo em Go:
- Quando você executa um programa Go, o sistema operacional cria um processo. Dentro desse processo, o runtime do Go é carregado e a função main é executada.

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    // Obtém o ID do processo (PID)
    pid := os.Getpid()
    fmt.Printf("Processo ID (PID): %d\n", pid)

    // Obtém o ID do processo pai (PPID)
    ppid := os.Getppid()
    fmt.Printf("Processo Pai ID (PPID): %d\n", ppid)

    // Simula algum trabalho
    time.Sleep(10 * time.Second)

    // Encerra o processo com um código de saída
    os.Exit(0)
}
```