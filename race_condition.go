package main

import (
	"fmt"
	"sync"
)

// O Problema:
// Queremos contar até 1000 usando 1000 goroutines.
// Cada goroutine adiciona +1 ao contador.
// Mas se você rodar, vai ver que o resultado NUNCA é 1000. Por que?

var count = 0

// Solução: Usamos um Mutex (Mutual Exclusion)
// É como uma chave de banheiro: só uma pessoa entra por vez.
var mu sync.Mutex

func increment() {
	mu.Lock() // 🔒 Tranca a porta
	count = count + 1
	mu.Unlock() // 🔓 Destranca a porta
}

func main() {
	// WaitGroup serve apenas para esperar todas as goroutines terminarem
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Contagem Final: %d (Esperado: 1000)\n", count)

	if count != 1000 {
		fmt.Println("❌ ERRO: Race Condition detectada!")
	} else {
		fmt.Println("✅ SUCESSO! O Mutex protegeu a variável.")
	}
}

// DESAFIO:
// Corrija este código para que a contagem seja sempre 1000.
// Dica: Você precisa impedir que duas goroutines mexam na variável 'count' ao mesmo tempo.
// Use 'sync.Mutex' para isso.
