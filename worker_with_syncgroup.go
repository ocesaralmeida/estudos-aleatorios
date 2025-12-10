// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// // Job representa uma tarefa a ser processada
// type Job struct {
// 	ID       int
// 	Duration time.Duration
// }

// // Result representa o resultado do processamento
// type Result struct {
// 	JobID int
// 	Err   error
// }

// // Worker: função que fará o processamento dos jobs
// // Ela recebe:
// // - id (número do worker, só pra log)
// // - jobs (canal SÓ PARA LEITURA)
// // - results (canal SÓ PARA ESCRITA)
// // - wg (WaitGroup para sinalizar que terminou sua execução)
// func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {

// 	// Quando o worker sair da função (acabou o range), sinalizamos ao WaitGroup
// 	// que este worker terminou.
// 	defer func() {
// 		fmt.Printf("[Worker %d] Finalizando worker e chamando wg.Done()\n", id)
// 		wg.Done()
// 	}()

// 	fmt.Printf("[Worker %d] Iniciado e aguardando jobs...\n", id)

// 	// Loop que lê jobs até que o canal seja fechado e esvaziado.
// 	for job := range jobs {
// 		fmt.Printf("[Worker %d] Recebeu job %d. Processando...\n", id, job.ID)

// 		// Simula um trabalho pesado
// 		time.Sleep(job.Duration)

// 		fmt.Printf("[Worker %d] Concluiu job %d. Enviando resultado...\n", id, job.ID)

// 		// Envia o resultado para o canal de resultados
// 		results <- Result{JobID: job.ID, Err: nil}
// 	}

// 	// Quando o canal jobs fecha, esse range termina automaticamente
// 	fmt.Printf("[Worker %d] Canal 'jobs' fechado. Nenhum job restante.\n", id)
// }

// func main() {
// 	const numJobs = 10
// 	const numWorkers = 3

// 	fmt.Println("==== Início da execução ====")

// 	// Criamos os dois canais bufferizados
// 	jobs := make(chan Job, numJobs)
// 	results := make(chan Result, numJobs)

// 	// WaitGroup vai esperar todos os workers finalizarem
// 	var wg sync.WaitGroup

// 	fmt.Printf("Iniciando %d workers...\n", numWorkers)

// 	// Adiciona o número de workers esperados no WaitGroup
// 	wg.Add(numWorkers)

// 	// Inicializa cada worker em sua própria goroutine
// 	for w := 1; w <= numWorkers; w++ {
// 		go worker(w, jobs, results, &wg)
// 	}

// 	// Envio dos jobs
// 	fmt.Printf("Enviando %d jobs ao canal...\n", numJobs)
// 	for j := 1; j <= numJobs; j++ {
// 		fmt.Printf("[Main] Enviando job %d\n", j)
// 		jobs <- Job{ID: j, Duration: 500 * time.Millisecond}
// 	}

// 	// Fechar o canal de jobs é ESSENCIAL para os workers saberem que não há mais trabalho
// 	fmt.Println("[Main] Fechando canal de jobs...")
// 	close(jobs)

// 	// Goroutine auxiliar para fechar o canal results assim que
// 	// todos os workers terminarem (quando wg.Wait() desbloquear)
// 	go func() {
// 		fmt.Println("[Goroutine Fechamento] Aguardando todos os workers finalizarem (wg.Wait)...")
// 		wg.Wait()
// 		fmt.Println("[Goroutine Fechamento] Todos os workers finalizaram. Fechando canal 'results'...")
// 		close(results)
// 	}()

// 	// Recebe resultados do canal results até ele ser fechado
// 	fmt.Println("[Main] Aguardando resultados...")
// 	for res := range results {
// 		fmt.Printf("[Main] Recebeu resultado do job %d\n", res.JobID)
// 	}

// 	fmt.Println("==== Todos os jobs processados! Encerrando programa. ====")
// }
