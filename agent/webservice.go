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

	consumerCtx, cancelConsumer := context.WithCancel(context.Background())
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second*5)

	go func() {
		log.Println("Consuming results from the Checker...")
		defer log.Println("Stopped Checker results consumption.")

		defer cancelConsumer()
		for r := range w.checkResultChan {
			w.checkResult = r
		}
	}()

	go func() {
		<-ctx.Done()
		defer cancelShutdown()

		log.Println("Closing pending HTTP connections...")

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("An error occurred while shutting down the agent web service: %v", err)
			return
		}

		log.Println("Pending connections closed.")
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-consumerCtx.Done()
	<-shutdownCtx.Done()

	return nil
}

func (w *webService) checkResultHandler(c *gin.Context) {
	c.JSON(200, w.checkResult)
}
