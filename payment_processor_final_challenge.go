// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"sync"
// 	"time"
// )

// // ==========================================
// // 1. DESIGN PATTERNS (Strategy, Factory, Decorator)
// // ==========================================

// // --- Strategy ---
// type PaymentStrategy interface {
// 	Process(amount float64) string
// }

// type Pix struct{}

// func (p *Pix) Process(amount float64) string {
// 	time.Sleep(100 * time.Millisecond) // Simula latência de rede
// 	return fmt.Sprintf("Pix de R$%.2f confirmado", amount)
// }

// type CreditCard struct{}

// func (c *CreditCard) Process(amount float64) string {
// 	time.Sleep(200 * time.Millisecond) // Crédito é mais lento
// 	return fmt.Sprintf("Crédito de R$%.2f aprovado", amount)
// }

// // --- Factory ---
// func NewPayment(method string) (PaymentStrategy, error) {
// 	switch method {
// 	case "pix":
// 		return &Pix{}, nil
// 	case "credit":
// 		return &CreditCard{}, nil
// 	default:
// 		return nil, fmt.Errorf("método inválido")
// 	}
// }

// // --- Decorator (Logging) ---
// type LogDecorator struct {
// 	Inner PaymentStrategy
// }

// func (l *LogDecorator) Process(amount float64) string {
// 	// Decorator adiciona comportamento antes
// 	// fmt.Printf("[LOG] Iniciando pagamento de %.2f...\n", amount)
// 	result := l.Inner.Process(amount)
// 	// Decorator adiciona comportamento depois
// 	return result + " (Logado)"
// }

// // ==========================================
// // 2. CONCORRÊNCIA (Mutex, WaitGroup, Channels)
// // ==========================================

// // --- Shared State (Mutex) ---
// // Queremos somar o total processado de forma segura.
// type SafeTotal struct {
// 	mu    sync.Mutex
// 	Value float64
// }

// func (s *SafeTotal) Add(amount float64) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.Value += amount
// }

// // --- Worker Pool ---
// type Job struct {
// 	ID     int
// 	Method string
// 	Amount float64
// }

// func worker(id int, jobs <-chan Job, total *SafeTotal, wg *sync.WaitGroup) {
// 	defer wg.Done() // Avisa que este worker terminou quando a função sair

// 	for job := range jobs {
// 		// 1. Factory: Cria a estratégia
// 		strategy, err := NewPayment(job.Method)
// 		if err != nil {
// 			fmt.Printf("Worker %d: Erro no job %d: %v\n", id, job.ID, err)
// 			continue
// 		}

// 		// 2. Decorator: Adiciona logs
// 		loggedStrategy := &LogDecorator{Inner: strategy}

// 		// 3. Strategy: Processa
// 		result := loggedStrategy.Process(job.Amount)
// 		fmt.Printf("Worker %d: %s\n", id, result)

// 		// 4. Mutex: Atualiza o total seguro
// 		total.Add(job.Amount)
// 	}
// }

// // ==========================================
// // MAIN (Orquestração)
// // ==========================================

// func main() {
// 	fmt.Println(">>> SISTEMA DE PAGAMENTOS MASSIVO <<<")
// 	start := time.Now()

// 	// Configuração
// 	const numJobs = 20
// 	const numWorkers = 5 // 5 workers processando 20 pagamentos

// 	jobs := make(chan Job, numJobs)
// 	var wg sync.WaitGroup
// 	totalProcessed := &SafeTotal{}

// 	// 1. Iniciar Workers (Goroutines)
// 	// O WaitGroup aqui serve para esperar os WORKERS terminarem, não os jobs individuais.
// 	// Mas neste padrão de worker pool, geralmente usamos WG para esperar o pool inteiro.
// 	for w := 1; w <= numWorkers; w++ {
// 		wg.Add(1)
// 		go worker(w, jobs, totalProcessed, &wg)
// 	}

// 	// 2. Enviar Jobs (Producer)
// 	methods := []string{"pix", "credit"}
// 	for j := 1; j <= numJobs; j++ {
// 		jobs <- Job{
// 			ID:     j,
// 			Method: methods[rand.Intn(len(methods))],
// 			Amount: float64(rand.Intn(100) + 10),
// 		}
// 	}
// 	close(jobs) // Fecha o canal para avisar os workers que não tem mais nada

// 	// 3. Esperar tudo terminar
// 	wg.Wait()

// 	duration := time.Since(start)
// 	fmt.Println("------------------------------------------------")
// 	fmt.Printf("✅ Processamento Concluído em %v\n", duration)
// 	fmt.Printf("💰 Total Processado (Seguro): R$%.2f\n", totalProcessed.Value)
// }
