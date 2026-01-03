// package main

// import (
// 	"errors"
// 	"fmt"
// )

// // ==========================================
// // PADRÃO 2: FACTORY
// // ==========================================
// // PROBLEMA: O código anterior era legal, mas quem cria o 'new(Pix)'?
// // Se o usuário manda uma string "pix" pelo frontend, precisamos de alguém
// // para converter essa string no objeto correto.

// // (Reaproveitando a Interface e Structs para o exemplo ficar completo)
// type PaymentStrategy interface {
// 	Process(amount float64)
// }

// type CreditCard struct{}

// func (c *CreditCard) Process(amount float64) { fmt.Printf("💳 Crédito: R$%.2f\n", amount) }

// type Pix struct{}

// func (p *Pix) Process(amount float64) { fmt.Printf("💠 Pix: R$%.2f\n", amount) }

// // A FÁBRICA (FACTORY)
// // A única responsabilidade dela é criar objetos.
// // Ela isola a complexidade de "escolher" qual objeto criar.
// func PaymentFactory(method string) (PaymentStrategy, error) {
// 	switch method {
// 	case "credito":
// 		return &CreditCard{}, nil
// 	case "pix":
// 		return &Pix{}, nil
// 	default:
// 		return nil, errors.New("método desconhecido")
// 	}
// }

// func main() {
// 	fmt.Println(">>> EXEMPLO 2: FACTORY <<<")

// 	// Simulando input do usuário (vinda de um JSON ou Frontend)
// 	inputs := []string{"pix", "credito", "boleto_invalido"}

// 	for _, input := range inputs {
// 		fmt.Printf("\nTentando pagar com: %s\n", input)

// 		// 1. Pedimos para a Fábrica criar a estratégia
// 		strategy, err := PaymentFactory(input)
// 		if err != nil {
// 			fmt.Printf("❌ Erro: %s\n", err)
// 			continue
// 		}

// 		// 2. Usamos a estratégia (sem saber qual é)
// 		strategy.Process(50.00)
// 	}
// }
