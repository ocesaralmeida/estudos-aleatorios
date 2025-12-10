// package main

// import (
// 	"fmt"
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

// // Worker é a função que processa os jobs
// func worker(id int, jobs <-chan Job, results chan<- Result) {
// 	for job := range jobs {
// 		fmt.Printf("Worker %d iniciou job %d\n", id, job.ID)
// 		time.Sleep(job.Duration)
// 		fmt.Printf("Worker %d finalizou job %d\n", id, job.ID)
// 		results <- Result{JobID: job.ID, Err: nil}
// 	}
// }

// func main() {
// 	const numJobs = 10
// 	const numWorkers = 3

// 	jobs := make(chan Job, numJobs)
// 	results := make(chan Result, numJobs)

// 	// 1. Iniciar Workers
// 	// Criamos 3 goroutines que ficam "ouvindo" o canal de jobs
// 	// Elas vão competir pelos jobs que chegarem
// 	for w := 1; w <= numWorkers; w++ {
// 		go worker(w, jobs, results)
// 	}

// 	// 2.Enviar Jobs
// 	for w := 1; w <= numJobs; w++ {
// 		jobs <- Job{ID: w, Duration: 500 * time.Millisecond}
// 	}
// 	close(jobs) // Importante: fechar o canal para avisar os workers que acabou

// 	// 3. Coletar Resultados
// 	// Como enviaros 10 jobs, precisamos esperar 10 respostas.
// 	// Se não fizermos isso, o programa (main) termina antes dos workers acabarem.
// 	for a := 1; a <= numJobs; a++ {
// 		res := <-results
// 		fmt.Printf("Main recebeu resultado do job %d\n", res.JobID)
// 	}

// 	fmt.Println("Todos os jobs processados")
// }
