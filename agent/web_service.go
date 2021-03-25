package agent

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type webService struct {
	checkResult     CheckResult
	checkResultChan <-chan CheckResult
	engine          *gin.Engine
}

func newWebService(checkResultChan <-chan CheckResult) *webService {
	w := &webService{
		checkResultChan: checkResultChan,
		engine:          gin.Default(),
	}

	w.engine.GET("/", w.checkResultHandler)

	return w
}

func (w *webService) Start(host string, port int, ctx context.Context) error {
	log.Println("Starting agent web service...")
	defer log.Println("Agent web service stopped.")

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", host, port),
		Handler:        w.engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("Consuming results from the Checker...")
		defer log.Println("Stopped Checker results consumption.")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				w.checkResult = <-w.checkResultChan
			}
		}
	}()

	pendingConnections := make(chan struct{})
	go func() {
		<-ctx.Done()

		log.Println("Closing pending HTTP connections...")
		defer log.Println("Pending connections closed.")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer shutdownCancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("An error occurred while shutting down the agent web service: %v", err)
		}
		close(pendingConnections)
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-pendingConnections

	return nil
}

func (w *webService) checkResultHandler(c *gin.Context) {
	c.JSON(200, w.checkResult)
}
