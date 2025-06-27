package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/middleware"
	"github.com/Zauro25/Capstone-PerpusKominfosan/routes"
	"github.com/Zauro25/Capstone-PerpusKominfosan/task"
	"time"
	"net/http"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	config.InitDB()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware dengan konfigurasi lebih ketat untuk production
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Gunakan env var, default "*" untuk development
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	
	if os.Getenv("ENV") == "development" {
		corsConfig.AllowOrigins = []string{"*"} // Lebih permisif di development
	}
	
	r.Use(cors.New(corsConfig))

	// Middleware tambahan
	r.Use(middleware.AuditLogMiddleware()) // Middleware untuk logging
	r.Use(gin.Recovery()) // Middleware recovery bawaan Gin

	// Jalankan scheduled tasks jika bukan environment test
	if os.Getenv("ENV") != "test" {
		tasks.StartNotificationTasks() // Perhatikan penulisan package 'task' bukan 'tasks'
	}

	// Setup routes
	routes.SetupRoutes(r)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    os.Getenv("ENV"),
		})
	})

	// Start server dengan graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		// Timeout config bisa ditambahkan di sini
	}

	log.Printf("Server starting on port %s in %s mode", port, os.Getenv("ENV"))
	
	// Jalankan server di goroutine untuk handle graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Tunggu interrupt signal untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Context dengan timeout untuk shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}