package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// ========== ESTRUTURAS DE DADOS ==========
type Invoice struct {
	ID     string
	Items  []Item
	Status string
	NFE    string // Nota Fiscal
}

type Item struct {
	ID       string
	Name     string
	Quantity int
	Price    float64
}

type Event struct {
	Type      string
	InvoiceID string
	Items     []Item
	Timestamp time.Time
}

// ========== CIRCUIT BREAKER ==========
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	state              CircuitState
	failures           int
	successes          int
	maxFailures        int
	resetTimeout       time.Duration
	lastFailure        time.Time
	mutex              sync.RWMutex
	halfOpenMaxAttempt int
	halfOpenAttempts   int
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:              StateClosed,
		maxFailures:        maxFailures,
		resetTimeout:       resetTimeout,
		halfOpenMaxAttempt: 3,
	}
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	if cb.state == StateClosed {
		return true
	}

	if cb.state == StateOpen {
		// Verifica se já pode tentar meio-aberto
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = StateHalfOpen
			cb.halfOpenAttempts = 0
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	}

	// StateHalfOpen: permite algumas tentativas
	if cb.halfOpenAttempts < cb.halfOpenMaxAttempt {
		return true
	}
	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateHalfOpen {
		cb.successes++
		if cb.successes >= cb.halfOpenMaxAttempt {
			cb.state = StateClosed
			cb.failures = 0
			cb.successes = 0
			log.Println("✅ Circuito FECHADO - Sefaz recuperado!")
		}
	} else if cb.state == StateClosed {
		cb.successes++
		if cb.successes > 10 { // Reset contadores após sucessos contínuos
			cb.failures = 0
			cb.successes = 0
		}
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.state == StateHalfOpen {
		cb.state = StateOpen // Volta para aberto se falhar no meio-aberto
		cb.halfOpenAttempts = 0
	} else if cb.state == StateClosed && cb.failures >= cb.maxFailures {
		cb.state = StateOpen
		log.Println("🔴 Circuito ABERTO - Sefaz instável!")
	}
}

func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// ========== SERVIÇOS EXTERNOS (MOCKS) ==========
type SefazClient struct {
	Name      string
	IsHealthy bool
	Cost      float64 // custo por chamada
}

func (s *SefazClient) GenerateInvoice(ctx context.Context, items []Item) (string, error) {
	// Simula latência
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	// Simula falha aleatória para Sefaz principal
	if s.Name == "Sefaz" && rand.Float32() < 0.3 { // 30% de falha
		s.IsHealthy = false
		return "", fmt.Errorf("sefaz: serviço indisponível")
	}

	s.IsHealthy = true
	nfe := fmt.Sprintf("NFE-%s-%d", s.Name, rand.Intn(10000))
	log.Printf("📄 %s gerou nota: %s (custo: R$%.2f)", s.Name, nfe, s.Cost)
	return nfe, nil
}

// ========== SEFAZ HANDLER ==========
type SefazHandler struct {
	circuitBreaker *CircuitBreaker
	sefazPrimary   *SefazClient
	sefazFallback  *SefazClient
	alertChan      chan<- Alert
	eventQueue     chan Event
	shutdown       chan struct{}
	wg             sync.WaitGroup
}

type Alert struct {
	Type    string
	Message string
	Level   string // "info", "warning", "error"
}

func NewSefazHandler() *SefazHandler {
	cb := NewCircuitBreaker(5, 10*time.Second) // 5 falhas, timeout 10s

	return &SefazHandler{
		circuitBreaker: cb,
		sefazPrimary: &SefazClient{
			Name: "Sefaz",
			Cost: 0.10, // R$0,10 por chamada
		},
		sefazFallback: &SefazClient{
			Name: "Sefaz_2",
			Cost: 1.50, // R$1,50 por chamada (mais caro!)
		},
		eventQueue: make(chan Event, 100),
		shutdown:   make(chan struct{}),
	}
}

func (h *SefazHandler) Start(ctx context.Context) {
	h.wg.Add(1)
	go h.processEvents(ctx)
	log.Println("🚀 SefazHandler iniciado")
}

func (h *SefazHandler) Stop() {
	close(h.shutdown)
	h.wg.Wait()
	log.Println("🛑 SefazHandler parado")
}

func (h *SefazHandler) SubmitEvent(event Event) {
	select {
	case h.eventQueue <- event:
		log.Printf("📨 Evento recebido: %s - Invoice: %s", event.Type, event.InvoiceID)
	default:
		log.Println("⚠️  Fila de eventos cheia, descartando evento")
	}
}

func (h *SefazHandler) processEvents(ctx context.Context) {
	defer h.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.shutdown:
			return
		case event := <-h.eventQueue:
			h.handleEvent(event)
		}
	}
}

func (h *SefazHandler) handleEvent(event Event) {
	invoiceID := event.InvoiceID

	// Tenta gerar nota fiscal
	nfe, err := h.generateNFE(event.Items)

	if err != nil {
		log.Printf("❌ Falha ao gerar NFE para invoice %s: %v", invoiceID, err)

		// Se falhar completamente, cancela compra
		h.cancelInvoice(invoiceID)
		return
	}

	// Sucesso - atualiza invoice
	h.updateInvoice(invoiceID, nfe)
	h.sendNotification(invoiceID, "success", "Compra realizada com sucesso!")
}

func (h *SefazHandler) generateNFE(items []Item) (string, error) {
	var nfe string
	var err error

	// Estratégia baseada no estado do Circuit Breaker
	state := h.circuitBreaker.GetState()

	switch state {
	case StateClosed:
		// Tenta Sefaz principal
		nfe, err = h.sefazPrimary.GenerateInvoice(context.Background(), items)
		if err != nil {
			h.circuitBreaker.RecordFailure()
			// Fallback imediato no primeiro erro
			return h.tryFallback(items)
		}
		h.circuitBreaker.RecordSuccess()
		return nfe, nil

	case StateOpen:
		// Usa fallback para maioria dos eventos
		// Mas 10% vai para teste (meio-aberto)
		if rand.Float32() < 0.1 { // 10% teste
			nfe, err = h.sefazPrimary.GenerateInvoice(context.Background(), items)
			if err == nil {
				h.circuitBreaker.RecordSuccess()
				log.Println("🎉 Sefaz recuperado durante teste!")
				return nfe, nil
			}
			h.circuitBreaker.RecordFailure()
		}

		// Usa fallback
		return h.tryFallback(items)

	case StateHalfOpen:
		// Modo teste - tenta principal
		nfe, err = h.sefazPrimary.GenerateInvoice(context.Background(), items)
		if err != nil {
			h.circuitBreaker.RecordFailure()
			return h.tryFallback(items)
		}
		h.circuitBreaker.RecordSuccess()
		return nfe, nil
	}

	return "", fmt.Errorf("estado desconhecido do circuit breaker")
}

func (h *SefazHandler) tryFallback(items []Item) (string, error) {
	// Gera alerta de uso do fallback
	h.sendAlert(Alert{
		Type:    "fallback_activated",
		Message: "Usando Sefaz_2 (serviço caro)",
		Level:   "warning",
	})

	// Tenta Sefaz_2
	nfe, err := h.sefazFallback.GenerateInvoice(context.Background(), items)
	if err != nil {
		// Falha total - ambos serviços indisponíveis
		h.sendAlert(Alert{
			Type:    "service_unavailable",
			Message: "Sefaz e Sefaz_2 indisponíveis!",
			Level:   "error",
		})
		return "", fmt.Errorf("todos os serviços de NFE indisponíveis")
	}

	return nfe, nil
}

func (h *SefazHandler) sendAlert(alert Alert) {
	log.Printf("🚨 ALERTA [%s]: %s", alert.Level, alert.Message)
	// Em produção, enviaria para sistema de monitoramento
}

func (h *SefazHandler) cancelInvoice(invoiceID string) {
	log.Printf("❌ Cancelando invoice %s", invoiceID)
	// 1. Cancela compra
	// 2. Gera extorno
	// 3. Envia notificação
	h.sendNotification(invoiceID, "cancelled", "Compra cancelada. Reembolso processado.")
}

func (h *SefazHandler) updateInvoice(invoiceID, nfe string) {
	log.Printf("✅ Invoice %s atualizado com NFE: %s", invoiceID, nfe)
	// Atualiza no banco de dados
}

func (h *SefazHandler) sendNotification(invoiceID, status, message string) {
	log.Printf("📢 Notificação: Invoice %s - %s: %s", invoiceID, status, message)
}

// ========== INVOICE SERVICE ==========
type InvoiceService struct {
	sefazHandler *SefazHandler
}

func (s *InvoiceService) ProcessPurchase(items []Item) string {
	invoiceID := fmt.Sprintf("INV-%d", rand.Intn(10000))

	invoice := Invoice{
		ID:     invoiceID,
		Items:  items,
		Status: "processing",
	}

	log.Printf("🛒 Processando compra: %s", invoiceID)

	// Publica evento para processamento assíncrono
	event := Event{
		Type:      "invoice.created",
		InvoiceID: invoiceID,
		Items:     items,
		Timestamp: time.Now(),
	}

	s.sefazHandler.SubmitEvent(event)

	return invoiceID
}

// ========== MAIN ==========
func main() {
	log.Println("🏪 Mercado Livre - Sistema de Notas Fiscais")
	log.Println("==========================================")

	// Inicializa serviços
	handler := NewSefazHandler()
	invoiceService := &InvoiceService{sefazHandler: handler}

	// Contexto para shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Inicia handler em background
	handler.Start(ctx)
	defer handler.Stop()

	// Simula algumas compras
	for i := 1; i <= 20; i++ {
		time.Sleep(500 * time.Millisecond)

		items := []Item{
			{ID: fmt.Sprintf("ITEM-%d", i), Name: "Produto Teste", Quantity: 1, Price: 99.90},
		}

		invoiceID := invoiceService.ProcessPurchase(items)
		log.Printf("📝 Compra %d iniciada: %s", i, invoiceID)

		// A cada 5 compras, mostra estado do circuit breaker
		if i%5 == 0 {
			state := handler.circuitBreaker.GetState()
			states := map[CircuitState]string{
				StateClosed:   "FECHADO",
				StateOpen:     "ABERTO",
				StateHalfOpen: "MEIO-ABERTO",
			}
			log.Printf("📊 Estado do Circuit Breaker: %s", states[state])
		}
	}

	// Aguarda processamento
	time.Sleep(3 * time.Second)
	log.Println("🎯 Simulação concluída!")
}
