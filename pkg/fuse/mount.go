package fuse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// MetadataFile disinkronkan dengan server, sekarang punya ModifiedTime
type MetadataFile struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"is_dir"`
	ModifiedTime int64  `json:"modified_time"` // Unix timestamp dari Host
}

// RemoteFileNode adalah node custom untuk file jarak jauh
type RemoteFileNode struct {
	fs.Inode
	metadata MetadataFile
}

// Getattr dipanggil sistem saat perintah `ls -la` atau `stat` dijalankan
func (n *RemoteFileNode) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0644                      // Izin rw-r--r--
	out.Size = uint64(n.metadata.Size)   // Ukuran asli dari Host
	
	// Set timestamp agar tidak muncul tahun 1970
	out.Mtime = uint64(n.metadata.ModifiedTime)
	out.Ctime = uint64(n.metadata.ModifiedTime)
	return 0
}

// NetUSBRoot adalah root folder virtual
type NetUSBRoot struct {
	fs.Inode
	FilesData []byte
}

func (r *NetUSBRoot) OnAdd(ctx context.Context) {
	var files []MetadataFile
	if err := json.Unmarshal(r.FilesData, &files); err != nil {
		log.Printf("Gagal mengurai JSON metadata: %v\n", err)
		return
	}

	for i, f := range files {
		if f.IsDir {
			continue // Skip folder sementara waktu
		}

		// Buat Inode dari struktur RemoteFileNode kita sendiri
		child := r.NewPersistentInode(
			ctx,
			&RemoteFileNode{metadata: f},
			fs.StableAttr{Ino: uint64(i + 2)},
		)
		r.AddChild(f.Name, child, false)
	}
}

// MountVirtualDrive tidak ada perubahan dari sebelumnya
func MountVirtualDrive(mountpoint string, filesData []byte) {
	fmt.Printf("Mounting FUSE to %s...\n", mountpoint)
	if err := os.MkdirAll(mountpoint, 0755); err != nil {
		log.Fatalf("Gagal membuat folder mount: %v", err)
	}

	server, err := fs.Mount(mountpoint, &NetUSBRoot{FilesData: filesData}, &fs.Options{
		MountOptions: fuse.MountOptions{
			Debug: false,
		},
	})
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	fmt.Printf("🔥 Success! Drive is logically mounted.\n")
	fmt.Printf("Buka terminal baru dan ketik: ls -la %s\n", mountpoint)
	server.Wait()
}