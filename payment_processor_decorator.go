// package main

// import "fmt"

// // ==========================================
// // PADRÃO 3: DECORATOR
// // ==========================================
// // PROBLEMA: Queremos adicionar Logs (ou métricas, ou validação) em todos os pagamentos.
// // Não queremos editar o código do Pix nem do Crédito para colocar "fmt.Println".

// type PaymentStrategy interface {
// 	Process(amount float64)
// }

// type Pix struct{}

// func (p *Pix) Process(amount float64) { fmt.Printf("💠 Pix: R$%.2f\n", amount) }

// // O DECORATOR
// // Ele "finge" ser uma estratégia (implementa a interface),
// // mas ele guarda uma estratégia real dentro dele.
// type LoggerDecorator struct {
// 	// O "recheio" (a estratégia real que será decorada)
// 	InnerStrategy PaymentStrategy
// }

// // Ele implementa o método Process
// func (l *LoggerDecorator) Process(amount float64) {
// 	fmt.Println("[LOG] Iniciando transação...") // Comportamento EXTRA (Antes)

// 	l.InnerStrategy.Process(amount) // Chama a original

// 	fmt.Println("[LOG] Transação finalizada.") // Comportamento EXTRA (Depois)
// }

// func main() {
// 	fmt.Println(">>> EXEMPLO 3: DECORATOR <<<")

// 	// 1. Criamos o objeto original (simples)
// 	pixSimples := &Pix{}

// 	fmt.Println("--- Sem Decorator ---")
// 	pixSimples.Process(100.00)

// 	fmt.Println("\n--- Com Decorator ---")
// 	// 2. Criamos o Decorator e colocamos o Pix dentro dele
// 	pixComLog := &LoggerDecorator{
// 		InnerStrategy: pixSimples,
// 	}

// 	// 3. Chamamos o Process do Decorator
// 	// Ele vai logar -> chamar o pix -> logar de novo.
// 	pixComLog.Process(100.00)
// }
