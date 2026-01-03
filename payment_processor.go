// package main

// import "fmt"

// // PaymentProcessor processa pagamentos de várias formas
// type PaymentProcessor struct{}

// func (p *PaymentProcessor) ProcessPayment(amount float64, method string) error {
// 	if method == "credit_card" {
// 		fmt.Printf("Processando pagamento de R$%.2f via Cartão de Crédito...\n", amount)
// 		// Lógica complexa de crédito...
// 		fmt.Println("Verificando limite...")
// 		fmt.Println("Taxa de 5% aplicada.")
// 		return nil
// 	} else if method == "debit_card" {
// 		fmt.Printf("Processando pagamento de R$%.2f via Cartão de Débito...\n", amount)
// 		// Lógica complexa de débito...
// 		fmt.Println("Verificando saldo...")
// 		fmt.Println("Sem taxas.")
// 		return nil
// 	} else if method == "pix" {
// 		fmt.Printf("Processando pagamento de R$%.2f via PIX...\n", amount)
// 		// Lógica complexa de PIX...
// 		fmt.Println("Gerando QR Code...")
// 		fmt.Println("Desconto de 10% aplicado!")
// 		return nil
// 	} else {
// 		return fmt.Errorf("método de pagamento inválido: %s", method)
// 	}
// }

// // func main() {
// // 	processor := &PaymentProcessor{}
// // 	processor.ProcessPayment(100.00, "credit_card")
// // 	processor.ProcessPayment(50.00, "pix")
// // }
