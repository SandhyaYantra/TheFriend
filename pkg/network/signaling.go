package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// MetadataFile menyimpan informasi file untuk dikirim lewat jaringan
type MetadataFile struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"is_dir"`
	ModifiedTime int64  `json:"modified_time"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StartSignalingServer sekarang menerima argumen sharedDir (direktori yang dibagikan)
func StartSignalingServer(port string, sharedDir string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to websocket:", err)
			return
		}
		defer conn.Close()

		fmt.Println("-> Klien terhubung! Meminta daftar file dari:", sharedDir)

		for {
			// Menunggu request (permintaan) dari klien
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Klien terputus:", err)
				break
			}
			
			// Membaca isi folder fisik di Host
			entries, err := os.ReadDir(sharedDir)
			if err != nil {
				log.Println("Gagal membaca direktori:", err)
				continue
			}

			// Mengumpulkan metadata dari setiap file/folder
			var files []MetadataFile
			for _, entry := range entries {
				info, err := entry.Info()
				if err != nil {
					continue
				}
				files = append(files, MetadataFile{
					Name:         entry.Name(),
					Size:         info.Size(),
					IsDir:        entry.IsDir(),
					ModifiedTime: info.ModTime().Unix(), // Mengambil waktu modifikasi asli
				})
			}

			// Marshal (Ubah) data Go menjadi format JSON
			jsonData, err := json.Marshal(files)
			if err != nil {
				log.Println("Gagal membuat JSON:", err)
				continue
			}

			// Kirim payload (muatan data) JSON kembali ke klien
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Println("Gagal mengirim data:", err)
				break
			}
		}
	})

	fmt.Printf("Signaling server is ALIVE on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server crashed: ", err)
	}
}