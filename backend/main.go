package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/spf13/cobra"

	"bitor/handlers"
	"bitor/middleware"
	_ "bitor/migrations" // ensure migrations are registered
	"bitor/nuclei"
	"bitor/routes"
	"bitor/scan"
	"bitor/setup"
	"bitor/ssh"
	"bitor/templates"
	"bitor/terminal"
	"bitor/utils"
	"bitor/utils/crypto"
	"bitor/version"
)

//go:embed all:pb_public
var distDir embed.FS
var distDirFS, _ = fs.Sub(distDir, "pb_public")

//go:embed ansible/roles/generate ansible/roles/terraform ansible/roles/nuclei ansible/generate.yml ansible/defaults.yml
var ansibleFiles embed.FS

var Version = "development"
var ansibleBasePath string
var showAnsibleLogs bool

// extractAnsibleFiles extracts embedded ansible files if they don't exist
func extractAnsibleFiles() error {
	// Check if ansible directory already exists
	if _, err := os.Stat(ansibleBasePath); err == nil {
		log.Printf("Ansible directory already exists at %s, skipping extraction", ansibleBasePath)
		return nil
	}

	log.Printf("Extracting ansible files to %s", ansibleBasePath)

	// Walk through the embedded files and extract them
	return fs.WalkDir(ansibleFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == "." {
			return nil
		}

		// Calculate relative path - strip the ansible/ prefix if it exists
		relPath := strings.TrimPrefix(path, "ansible/")

		// Skip if relPath is empty after trimming
		if relPath == "" {
			return nil
		}

		// Create target path
		targetPath := filepath.Join(ansibleBasePath, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetPath, err)
			}
			return nil
		}

		// Read the file content
		content, err := ansibleFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %v", path, err)
		}

		// Create parent directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %v", targetPath, err)
		}

		// Write the file
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", targetPath, err)
		}

		log.Printf("Extracted %s", targetPath)
		return nil
	})
}

func init() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
		cwd = "."
	}

	// Always use ansible directory in current working directory
	ansibleBasePath = filepath.Join(cwd, "ansible")
	nuclei.SetTemplatesDir(filepath.Join(cwd, "nuclei-templates"))
	log.Printf("Using base path: %s", cwd)

	setup.Version = Version
	version.Version = Version // Set the version in the version package
}

func main() {
	// Get the absolute path of the working directory
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Ensure we have a data directory
	dataDir := filepath.Join(workDir, "pb_data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	app := pocketbase.New()

	// Set the data directory and bind address explicitly
	//app.RootCmd.SetArgs([]string{"serve", "--dir", dataDir, "0.0.0.0:8090"})

	// Extract ansible files if needed
	if err := extractAnsibleFiles(); err != nil {
		log.Fatal(err)
	}

	// Initialize nuclei templates if needed
	if err := nuclei.InitializeTemplates(); err != nil {
		log.Fatal(err)
	}

	// Initialize templates for file browser
	if err := templates.InitializeTemplates(); err != nil {
		log.Fatal(err)
	}

	// Register migrations with automigrate enabled
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	// Initialize database and run migrations
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	// Add command line flags
	app.RootCmd.PersistentFlags().BoolVar(&showAnsibleLogs, "show-ansible-logs", false, "Show Ansible logs in the terminal")

	// Set environment variable for ansible logs based on flag
	app.RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if showAnsibleLogs {
			os.Setenv("SHOW_ANSIBLE_LOGS", "true")
			log.Println("Ansible logs will be shown in terminal")
		}
	}

	// Only register the ansible base path flag if it's not already set via environment variable
	if os.Getenv("ANSIBLE_BASE_PATH") == "" {
		log.Printf("Ansible base path before flag: %s", ansibleBasePath)
		app.RootCmd.PersistentFlags().StringVar(
			&ansibleBasePath,
			"ansible-base-path",
			ansibleBasePath,
			"base path for ansible playbooks",
		)
		log.Printf("Ansible base path after flag registration: %s", ansibleBasePath)
	} else {
		log.Printf("Using ansible base path from environment: %s", ansibleBasePath)
	}

	// Configure file serving and services after migrations are complete
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		log.Printf("Ansible base path in OnBeforeServe: %s", ansibleBasePath)

		// Add request logging middleware (before other middlewares)
		e.Router.Use(middleware.RequestLogger())
		log.Printf("Request logging middleware enabled")

		// Initialize the app
		if err := setup.InitializeApp(app); err != nil {
			return err
		}

		// Run setup after migrations
		if err := setup.Setup(app); err != nil {
			return err
		}

		// Set ansible base path for SSH operations and initialize keys if needed
		ssh.SetAnsibleBasePath(ansibleBasePath)
		if err := ssh.InitializeSSHKeys(app); err != nil {
			log.Printf("Failed to initialize SSH keys: %v", err)
			return err
		}

		// Verify SSH keys match between database and files
		if err := ssh.VerifySSHKeys(app); err != nil {
			log.Printf("Failed to verify SSH keys: %v", err)
			return err
		}

		// Ensure default groups exist
		if err := setup.EnsureGroupsCollection(app); err != nil {
			return err
		}

		// Validate encryption key after migrations
		if err := crypto.ValidateEncryptionKey(app); err != nil {
			return err
		}

		// Initialize notification service
		notificationService, err := routes.InitNotificationService(app)
		if err != nil {
			log.Printf("Failed to initialize notification service: %v", err)
			return err
		}

		// Initialize scan handlers
		scan.InitHandlers(app, ansibleBasePath, notificationService)

		// Terminal WebSocket handler
		log.Printf("Registering terminal route...")
		e.Router.GET("/api/terminal", terminal.HandleTerminalConnection(app))

		// Register version check routes
		version.RegisterRoutes(e)

		// Register routes and services
		if err := routes.RegisterRoutes(app, ansibleBasePath, notificationService, e); err != nil {
			log.Printf("Failed to register routes: %v", err)
			return err
		}

		// Expose debug routes for development purposes
		handlers.RegisterDebugRoutes(e)

		// Initialize attack surface handlers and register routes directly
		attackSurfaceHandlers := handlers.NewAttackSurfaceHandlers(app)

		// Register attack surface routes directly here
		fmt.Println("DEBUG: Registering attack surface routes directly in main OnBeforeServe")
		e.Router.GET("/api/attack-surface/sources", func(c echo.Context) error {
			fmt.Println("DEBUG: Sources endpoint called in main!")
			return attackSurfaceHandlers.HandleGetAvailableSources(c)
		})
		e.Router.GET("/api/attack-surface/subdomains", func(c echo.Context) error {
			fmt.Println("DEBUG: Subdomains endpoint called in main!")
			return attackSurfaceHandlers.HandleGetSubdomains(c)
		})
		e.Router.GET("/api/attack-surface/subdomains/stats", func(c echo.Context) error {
			fmt.Println("DEBUG: Stats endpoint called in main!")
			return attackSurfaceHandlers.HandleGetDomainStats(c)
		})
		e.Router.POST("/api/attack-surface/subdomains/scan", func(c echo.Context) error {
			fmt.Println("DEBUG: Scan endpoint called in main!")
			return attackSurfaceHandlers.HandleStartSubdomainScan(c)
		})
		e.Router.POST("/api/attack-surface/tld/discover", func(c echo.Context) error {
			fmt.Println("DEBUG: TLD Discovery endpoint called in main!")
			return attackSurfaceHandlers.HandleStartTLDDiscovery(c)
		})
		e.Router.GET("/api/attack-surface/tld", func(c echo.Context) error {
			fmt.Println("DEBUG: TLD Get endpoint called in main!")
			return attackSurfaceHandlers.HandleGetTLDs(c)
		})

		// Netblock discovery routes
		e.Router.POST("/api/attack-surface/netblock/discover", func(c echo.Context) error {
			fmt.Println("DEBUG: Netblock Discovery endpoint called in main!")
			return attackSurfaceHandlers.HandleStartNetblockDiscovery(c)
		})
		e.Router.GET("/api/attack-surface/netblocks", func(c echo.Context) error {
			fmt.Println("DEBUG: Netblocks Get endpoint called in main!")
			return attackSurfaceHandlers.HandleGetNetblocks(c)
		})
		e.Router.GET("/api/attack-surface/ips", func(c echo.Context) error {
			fmt.Println("DEBUG: IPs Get endpoint called in main!")
			return attackSurfaceHandlers.HandleGetIPs(c)
		})
		e.Router.GET("/api/attack-surface/netblocks/stats", func(c echo.Context) error {
			fmt.Println("DEBUG: Netblock Stats endpoint called in main!")
			return attackSurfaceHandlers.HandleGetNetblockStats(c)
		})

		// Port scanning routes
		e.Router.POST("/api/attack-surface/ports/scan", func(c echo.Context) error {
			fmt.Println("DEBUG: Port Scan endpoint called in main!")
			return attackSurfaceHandlers.HandleStartPortScan(c)
		})
		e.Router.GET("/api/attack-surface/ports/scan/:scan_id/progress", func(c echo.Context) error {
			fmt.Println("DEBUG: Port Scan Progress endpoint called in main!")
			return attackSurfaceHandlers.HandleGetPortScanProgress(c)
		})
		e.Router.GET("/api/attack-surface/ports", func(c echo.Context) error {
			fmt.Println("DEBUG: Ports Get endpoint called in main!")
			return attackSurfaceHandlers.HandleGetPorts(c)
		})
		e.Router.GET("/api/attack-surface/ports/scans", func(c echo.Context) error {
			fmt.Println("DEBUG: Port Scans Get endpoint called in main!")
			return attackSurfaceHandlers.HandleGetPortScans(c)
		})
		e.Router.GET("/api/attack-surface/ports/stats", func(c echo.Context) error {
			fmt.Println("DEBUG: Port Stats endpoint called in main!")
			return attackSurfaceHandlers.HandleGetPortStats(c)
		})
		e.Router.POST("/api/attack-surface/urls/scan", func(c echo.Context) error {
			fmt.Println("DEBUG: URL Scan endpoint called in main!")
			return attackSurfaceHandlers.HandleStartURLScan(c)
		})
		e.Router.GET("/api/attack-surface/urls/scan/:scan_id/progress", func(c echo.Context) error {
			fmt.Println("DEBUG: URL Scan Progress endpoint called in main!")
			return attackSurfaceHandlers.HandleGetURLScanProgress(c)
		})
		e.Router.GET("/api/attack-surface/urls", func(c echo.Context) error {
			fmt.Println("DEBUG: URLs Get endpoint called in main!")
			return attackSurfaceHandlers.HandleGetURLs(c)
		})
		e.Router.GET("/api/attack-surface/urls/scans", func(c echo.Context) error {
			fmt.Println("DEBUG: URL Scans Get endpoint called in main!")
			return attackSurfaceHandlers.HandleGetURLScans(c)
		})
		e.Router.GET("/api/attack-surface/urls/stats", func(c echo.Context) error {
			fmt.Println("DEBUG: URL Stats endpoint called in main!")
			return attackSurfaceHandlers.HandleGetURLStats(c)
		})
		e.Router.POST("/api/attack-surface/nuclei/collect-targets", func(c echo.Context) error {
			fmt.Println("DEBUG: Nuclei Collect Targets endpoint called in main!")
			return attackSurfaceHandlers.HandleCollectAttackSurfaceTargets(c)
		})
		e.Router.POST("/api/attack-surface/nuclei/create-target", func(c echo.Context) error {
			fmt.Println("DEBUG: Nuclei Create Target endpoint called in main!")
			return attackSurfaceHandlers.HandleCreateNucleiTargetFromAttackSurface(c)
		})

		fmt.Println("DEBUG: Attack surface routes registered successfully in main!")
		fmt.Println("DEBUG: GET /api/attack-surface/sources -> HandleGetAvailableSources")
		fmt.Println("DEBUG: GET /api/attack-surface/subdomains -> HandleGetSubdomains")
		fmt.Println("DEBUG: GET /api/attack-surface/subdomains/stats -> HandleGetDomainStats")
		fmt.Println("DEBUG: POST /api/attack-surface/subdomains/scan -> HandleStartSubdomainScan")
		fmt.Println("DEBUG: POST /api/attack-surface/tld/discover -> HandleStartTLDDiscovery")
		fmt.Println("DEBUG: GET /api/attack-surface/tld -> HandleGetTLDs")
		fmt.Println("DEBUG: POST /api/attack-surface/netblock/discover -> HandleStartNetblockDiscovery")
		fmt.Println("DEBUG: GET /api/attack-surface/netblocks -> HandleGetNetblocks")
		fmt.Println("DEBUG: GET /api/attack-surface/ips -> HandleGetIPs")
		fmt.Println("DEBUG: GET /api/attack-surface/netblocks/stats -> HandleGetNetblockStats")
		fmt.Println("DEBUG: POST /api/attack-surface/ports/scan -> HandleStartPortScan")
		fmt.Println("DEBUG: GET /api/attack-surface/ports -> HandleGetPorts")
		fmt.Println("DEBUG: GET /api/attack-surface/ports/scans -> HandleGetPortScans")
		fmt.Println("DEBUG: GET /api/attack-surface/ports/stats -> HandleGetPortStats")
		fmt.Println("DEBUG: POST /api/attack-surface/urls/scan -> HandleStartURLScan")
		fmt.Println("DEBUG: GET /api/attack-surface/urls -> HandleGetURLs")
		fmt.Println("DEBUG: GET /api/attack-surface/urls/scans -> HandleGetURLScans")
		fmt.Println("DEBUG: GET /api/attack-surface/urls/stats -> HandleGetURLStats")
		fmt.Println("DEBUG: POST /api/attack-surface/nuclei/collect-targets -> HandleCollectAttackSurfaceTargets")
		fmt.Println("DEBUG: POST /api/attack-surface/nuclei/create-target -> HandleCreateNucleiTargetFromAttackSurface")

		// Serve static files from pb_public directory
		e.Router.GET("/*", echo.WrapHandler(http.FileServer(http.FS(distDirFS))))

		return nil
	})

	// Add password change middleware
	setup.AddPasswordChangeMiddleware(app)

	// Check for required dependencies
	if err := utils.CheckDependencies(); err != nil {
		log.Fatal(err)
	}

	// Start the application
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
// Railway test deployment Fri Jul 11 13:43:19 BST 2025
