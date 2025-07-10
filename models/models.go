package models

import (
	"time"
	"gorm.io/gorm"
)

// Base model with common fields
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Admin Perpustakaan model
type AdminPerpustakaan struct {
	BaseModel
	Username       string       `json:"username" gorm:"uniqueIndex;not null"`
	Password       string       `json:"-" gorm:"not null"`
	NamaLengkap    string       `json:"nama_lengkap" gorm:"not null"`
	Email          string       `json:"email" gorm:"uniqueIndex"`
	NoTelepon      string       `json:"no_telepon"`
	PerpustakaanID uint         `json:"perpustakaan_id" gorm:"not null"`
	Perpustakaan   Perpustakaan `json:"perpustakaan" gorm:"foreignKey:PerpustakaanID"`
	IsActive       bool         `json:"is_active" gorm:"default:false"`
	LastLogin      *time.Time   `json:"last_login"`
}

// Admin DPK model
type AdminDPK struct {
	BaseModel
	Username    string     `json:"username" gorm:"uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"not null"`
	NamaLengkap string     `json:"nama_lengkap" gorm:"not null"`
	Email       string     `json:"email" gorm:"uniqueIndex"`
	NoTelepon   string     `json:"no_telepon"`
	HakAkses    string     `json:"hak_akses" gorm:"default:'admin'"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	LastLogin   *time.Time `json:"last_login"`
}

// Executive model
type Executive struct {
	BaseModel
	Username    string     `json:"username" gorm:"uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"not null"`
	NamaLengkap string     `json:"nama_lengkap" gorm:"not null"`
	Email       string     `json:"email" gorm:"uniqueIndex"`
	NoTelepon   string    `json:"no_telepon"`
	Jabatan     string     `json:"jabatan"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	LastLogin   *time.Time `json:"last_login"`
}

// Perpustakaan model
type Perpustakaan struct {
	BaseModel
	Periode 		string     `json:"periode" gorm:"not null"` 
	NamaPerpustakaan string     `json:"nama_perpustakaan" gorm:"not null"`
	Alamat           string     `json:"alamat" gorm:"not null"`
	KepalaPerpustakaan string `json:"kepala_perpustakaan"` // Nama kepala perpustakaan
	JenisPerpustakaan string    `json:"jenis_perpustakaan" gorm:"not null"` // Umum/Sekolah/Khusus/Perguruan Tinggi
	TahunBerdiri	 int `json:"tahun_berdiri"` // Tahun berdiri perpustakaan
	NomorInduk       int     `json:"nomor_induk" gorm:"uniqueIndex"`
	JumlahSDM	   int        `json:"jumlah_sdm" gorm:"default:0"` // Jumlah SDM yang terdaftar
	JumlahPengunjung int    `json:"jumlah_pengunjung" gorm:"default:0"` // Jumlah pengunjung yang terdaftar
	JumlahAnggota	int        `json:"jumlah_anggota" gorm:"default:0"` // Jumlah anggota yang terdaftar
	StatusVerifikasi string     `json:"status_verifikasi" gorm:"default:'Pending'"` // Pending/Terkirim/Disetujui/Direvisi
	TanggalKirim     *time.Time `json:"tanggal_kirim"`
	CatatanRevisi    string     `json:"catatan_revisi"`
	CreatedBy        uint       `json:"created_by" gorm:"index"` // User ID pembuat data
	// Relations
	AdminPerpustakaan []AdminPerpustakaan `json:"admin_perpustakaan" gorm:"foreignKey:PerpustakaanID"`
	Koleksi           []Koleksi           `json:"koleksi" gorm:"foreignKey:PerpustakaanID"`
	SDM               []SDM               `json:"sdm" gorm:"foreignKey:PerpustakaanID"`
	Pengunjung        []Pengunjung        `json:"pengunjung" gorm:"foreignKey:PerpustakaanID"`
	Anggota           []Anggota           `json:"anggota" gorm:"foreignKey:PerpustakaanID"`
}

// Koleksi model
type Koleksi struct {
	BaseModel
	PerpustakaanID   uint   `json:"perpustakaan_id" gorm:"not null"`
	JenisKoleksi     string `json:"jenis_koleksi" gorm:"not null"` // Buku/Majalah/Jurnal/dll
	JumlahKoleksi    int    `json:"jumlah_koleksi" gorm:"not null"`
	KondisiKoleksi   string `json:"kondisi_koleksi"` // Baik/Rusak/Hilang
	TanggalPerolehan *time.Time `json:"tanggal_perolehan"`
	Keterangan       string `json:"keterangan"`
}

// SDM model
type SDM struct {
	BaseModel
	PerpustakaanID uint   `json:"perpustakaan_id" gorm:"not null"`
	JenisSDM       string `json:"jenis_sdm" gorm:"not null"` // Pustakawan/Staf/Tenaga Teknis
	JumlahSDM      int    `json:"jumlah_sdm" gorm:"not null"`
	Pendidikan     string `json:"pendidikan"`
	Status         string `json:"status"` // PNS/Honorer/Kontrak
	Keterangan     string `json:"keterangan"`
}

// Pengunjung model
type Pengunjung struct {
	BaseModel
	PerpustakaanID    uint      `json:"perpustakaan_id" gorm:"not null"`
	TanggalKunjungan  time.Time `json:"tanggal_kunjungan" gorm:"not null"`
	JumlahPengunjung  int       `json:"jumlah_pengunjung" gorm:"not null"`
	JenisKelamin      string    `json:"jenis_kelamin"` // L/P
	KelompokUsia      string    `json:"kelompok_usia"` // Anak/Remaja/Dewasa/Lansia
	TujuanKunjungan   string    `json:"tujuan_kunjungan"`
	Keterangan        string    `json:"keterangan"`
}

// Anggota model
type Anggota struct {
	BaseModel
	PerpustakaanID   uint      `json:"perpustakaan_id" gorm:"not null"`
	NamaAnggota      string    `json:"nama_anggota" gorm:"not null"`
	NomorAnggota     string    `json:"nomor_anggota" gorm:"uniqueIndex;not null"`
	JenisKelamin     string    `json:"jenis_kelamin"` // L/P
	TanggalLahir     *time.Time `json:"tanggal_lahir"`
	Alamat           string    `json:"alamat"`
	NoTelepon        string    `json:"no_telepon"`
	TanggalDaftar    time.Time `json:"tanggal_daftar" gorm:"not null"`
	StatusAnggota    string    `json:"status_anggota" gorm:"default:'Aktif'"` // Aktif/Tidak Aktif
	Keterangan       string    `json:"keterangan"`
}

// Verifikasi model
type Verifikasi struct {
	BaseModel
	PerpustakaanID uint         `json:"perpustakaan_id" gorm:"not null"`
	Perpustakaan   Perpustakaan `json:"perpustakaan" gorm:"foreignKey:PerpustakaanID"`
	JenisData      string       `json:"jenis_data" gorm:"not null"` // Perpustakaan/Koleksi/SDM/dll
	Status         string       `json:"status" gorm:"not null"` // Pending/Disetujui/Direvisi
	CatatanRevisi  string       `json:"catatan_revisi"`
	TanggalVerifikasi *time.Time `json:"tanggal_verifikasi"`
	AdminDPKID     uint         `json:"admin_dpk_id" gorm:"default:null"`
	AdminDPK       AdminDPK     `json:"admin_dpk" gorm:"foreignKey:AdminDPKID"`
}

// Laporan model
type Laporan struct {
	BaseModel
	Judul          string    `json:"judul" gorm:"not null"`
	Periode        string    `json:"periode" gorm:"not null"` // 2024-1, 2024-2
	JenisLaporan   string    `json:"jenis_laporan" gorm:"not null"` // Statistik/Rekapitulasi/Evaluasi
	FilePath       string    `json:"file_path"`
	FormatFile     string    `json:"format_file"` // PDF/Excel/CSV
	TanggalGenerate time.Time `json:"tanggal_generate" gorm:"not null"`
	AdminDPKID     uint      `json:"admin_dpk_id" gorm:"not null"`
	AdminDPK       AdminDPK  `json:"admin_dpk" gorm:"foreignKey:AdminDPKID"`
	Status         string    `json:"status" gorm:"default:'Generated'"` // Generated/Sent/Archived
	ChartData	   string    `json:"chart_data" gorm:"type:text"` // JSON untuk chart (optional)
}

// LogRevisi model
type LogRevisi struct {
	BaseModel
	PerpustakaanID uint         `json:"perpustakaan_id" gorm:"not null"`
	Perpustakaan   Perpustakaan `json:"perpustakaan" gorm:"foreignKey:PerpustakaanID"`
	JenisRevisi    string       `json:"jenis_revisi" gorm:"not null"`
	CatatanRevisi  string       `json:"catatan_revisi" gorm:"not null"`
	StatusSebelum  string       `json:"status_sebelum" gorm:"not null"`
	StatusSesudah  string       `json:"status_sesudah" gorm:"not null"`
	AdminDPKID     uint         `json:"admin_dpk_id"`
	AdminDPK       AdminDPK     `json:"admin_dpk" gorm:"foreignKey:AdminDPKID"`
	TanggalRevisi  time.Time    `json:"tanggal_revisi" gorm:"not null"`
}

// Notifikasi model
type Notifikasi struct {
       BaseModel
       JudulNotifikasi string    `json:"judul_notifikasi" gorm:"not null"`
       IsiNotifikasi   string    `json:"isi_notifikasi" gorm:"not null;type:text"`
       JenisNotifikasi string    `json:"jenis_notifikasi" gorm:"not null;type:varchar(20)"`
       TujuanUser      string    `json:"tujuan_user" gorm:"not null;type:varchar(50)"`
       UserID          *uint     `json:"user_id" gorm:"index"` // Null untuk broadcast
       RelatedID       *uint     `json:"related_id"` // ID terkait (perpustakaan, verifikasi, dll)
       RelatedType     string    `json:"related_type"` // Jenis relasi
       IsRead          bool      `json:"is_read" gorm:"default:false"`
       TanggalKirim    time.Time `json:"tanggal_kirim" gorm:"not null"`
       ExpiredAt       time.Time `json:"expired_at"` // Waktu kadaluarsa notifikasi
       ActionLink      string    `json:"action_link"` // Link aksi (optional)
}

// AuditLog model
type AuditLog struct {
	BaseModel
	UserType   string    `json:"user_type" gorm:"not null"` // admin_perpustakaan/admin_dpk/executive
	UserID     uint      `json:"user_id" gorm:"not null"`
	Action     string    `json:"action" gorm:"not null"` // Create/Update/Delete/Login/Logout
	TableName  string    `json:"table_name"`
	RecordID   uint      `json:"record_id"`
	OldValues  string    `json:"old_values" gorm:"type:text"`
	NewValues  string    `json:"new_values" gorm:"type:text"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Timestamp  time.Time `json:"timestamp" gorm:"not null"`
}

// DTOs for API requests/responses
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"user_type"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      interface{} `json:"user"`
	UserType  string      `json:"user_type"`
	ExpiresAt int64       `json:"expires_at"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type SendDataRequest struct {
	PerpustakaanID uint   `json:"perpustakaan_id" binding:"required"`
}

type VerifikasiRequest struct {
	PerpustakaanID uint   `json:"perpustakaan_id" binding:"required"`
	Status         string `json:"status" binding:"required"` // Disetujui/Direvisi
	CatatanRevisi  string `json:"catatan_revisi"`
}

type LaporanRequest struct {
    Judul        string `json:"judul" binding:"required"`
    Periode      string `json:"periode" binding:"required"`
    JenisLaporan string `json:"jenis_laporan" binding:"required"`
    FormatFile   string `json:"format_file" binding:"required"`
}

type NotifikasiRequest struct {
	JudulNotifikasi string `json:"judul_notifikasi" binding:"required"`
	IsiNotifikasi   string `json:"isi_notifikasi" binding:"required"`
	JenisNotifikasi string `json:"jenis_notifikasi" binding:"required"`
	TujuanUser      string `json:"tujuan_user" binding:"required"`
	UserID          *uint  `json:"user_id"` // Optional, jika kosong berarti broadcast
}

type RegisterRequest struct {
	Username          string `json:"username" binding:"required"`
	Password          string `json:"password" binding:"required,min=6"`
	NamaLengkap       string `json:"nama_lengkap" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	NoTelepon         string `json:"no_telepon"`
	UserType          string `json:"user_type" binding:"required"` // admin_perpustakaan only
	NamaPerpustakaan  string `json:"nama_perpustakaan" binding:"required"`
	Alamat            string `json:"alamat" binding:"required"`
	JenisPerpustakaan string `json:"jenis_perpustakaan" binding:"required"` // Umum/Sekolah/Khusus
	NomorInduk        string `json:"nomor_induk" binding:"required"`
}

type RegisterAdminDPKRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	NoTelepon   string `json:"no_telepon"`
	UserType    string `json:"user_type" binding:"required"` // admin_dpk only
}
type VerifyAdminRequest struct {
    AdminPerpustakaanID uint   `json:"admin_perpustakaan_id" binding:"required"`
    Status              string `json:"status" binding:"required,oneof=approved rejected"`
    Catatan            string `json:"catatan"`
}
type InputDataPerpustakaanRequest struct {
	Periode		  string `json:"periode" binding:"required"` // Format YYYY-1 atau YYYY-2
    NamaPerpustakaan  string `json:"nama_perpustakaan" binding:"required"`
    Alamat            string `json:"alamat" binding:"required"`
    JenisPerpustakaan string `json:"jenis_perpustakaan" binding:"required,oneof=Umum Sekolah Khusus"`
	KepalaPerpustakaan string `json:"kepala_perpustakaan"`
	TahunBerdiri     int    `json:"tahun_berdiri" binding:"required,min=1900,max=2100"` // Validasi tahun
    NomorInduk        int `json:"nomor_induk" binding:"required"`
    JumlahSDM         int    `json:"jumlah_sdm"`
    JumlahPengunjung  int    `json:"jumlah_pengunjung"`
    JumlahAnggota     int    `json:"jumlah_anggota"`
	TanggalKirim      *time.Time `json:"tanggal_kirim"` 
}