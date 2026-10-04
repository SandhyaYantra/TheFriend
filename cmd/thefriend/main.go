package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SandhyaYantra/thefriend/pkg/fuse"
	"github.com/SandhyaYantra/thefriend/pkg/network"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "thefriend",
	Short: "TheFriend is a CLI tool to mount remote devices locally",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the daemon on the host device",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		
		fmt.Printf("Starting TheFriend Host Daemon...\n")
		fmt.Printf("Sharing directory: %s\n", dir)
		
		port := "8080"
		fmt.Println("Waiting for connection... (Mockup Token: X7Y9-B2V1)")
		
		// Server memblokir di baris ini
		network.StartSignalingServer(port, dir)
	},
}

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount a remote drive on the client device",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		mountpoint, _ := cmd.Flags().GetString("mountpoint")
		
		fmt.Printf("Connecting to remote host with token: %s\n", token)
		
		// Target URL ke server lokal
		url := "ws://localhost:8080/ws"
		fmt.Printf("Dialing signaling server at %s...\n", url)
		
		// Inisiasi koneksi WebSocket sebagai klien
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Fatal("Connection failed: ", err)
		}
		defer conn.Close()

		// Mengirim pesan "LIST" (Minta daftar file) ke server
		message := []byte("LIST")
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Fatal("Error sending message: ", err)
		}

		// Membaca balasan berupa JSON dari server
		_, responseJSON, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading response: ", err)
		}

		fmt.Printf("Berhasil menerima metadata (%d bytes) dari server.\n", len(responseJSON))
		
		// Eksekusi pemasangan virtual drive dengan membawa data JSON
		fuse.MountVirtualDrive(mountpoint, responseJSON)
	},
}

func init() {
	serveCmd.Flags().StringP("dir", "d", ".", "Directory to share")
	mountCmd.Flags().StringP("token", "t", "", "Pairing token from the host")
	mountCmd.Flags().StringP("mountpoint", "m", "./remote_usb", "Local path to mount the drive")
	mountCmd.MarkFlagRequired("token")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(mountCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}