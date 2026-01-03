// package main

// import "fmt"

// // ==========================================
// // PADRÃO 1: STRATEGY
// // ==========================================
// // PROBLEMA: Temos vários algoritmos para fazer a mesma coisa (pagar),
// // e queremos alternar entre eles facilmente sem encher o código de if/else.

// // 1. A Interface (O Contrato)
// // Todo mundo que quiser ser uma forma de pagamento TEM que ter esse método.
// type PaymentStrategy interface {
// 	Process(amount float64)
// }

// // 2. As Estratégias (As Implementações)

// // Cartão de Crédito
// type CreditCard struct{}

// func (c *CreditCard) Process(amount float64) {
// 	fmt.Printf("💳 Pagando R$%.2f com Crédito (Taxa 5%%)\n", amount)
// }

// // Pix
// type Pix struct{}

// func (p *Pix) Process(amount float64) {
// 	fmt.Printf("💠 Pagando R$%.2f com PIX (Desconto 10%%)\n", amount)
// }

// // 3. O Contexto (Quem usa)
// // Essa função aceita QUALQUER coisa que respeite a interface PaymentStrategy.
// // Ela não sabe se é Pix ou Crédito, e não importa!
// func PayOrder(amount float64, strategy PaymentStrategy) {
// 	fmt.Println("--- Iniciando Pedido ---")
// 	strategy.Process(amount)
// 	fmt.Println("--- Pedido Finalizado ---")
// }

// func main() {
// 	fmt.Println(">>> EXEMPLO 1: STRATEGY <<<")

// 	valor := 100.00

// 	// Cenário A: Usuário escolheu Crédito
// 	// Instanciamos a estratégia de crédito e passamos para a função.
// 	credito := &CreditCard{}
// 	PayOrder(valor, credito)

// 	fmt.Println()

// 	// Cenário B: Usuário escolheu Pix
// 	// Instanciamos a estratégia de Pix. Note que a função PayOrder é a mesma!
// 	pix := &Pix{}
// 	PayOrder(valor, pix)
// }
