package controllers

import (
	"net/http"
	"time"
	"encoding/json"
	"strconv"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/middleware"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
	
)

func Login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    now := time.Now()
    var response models.LoginResponse
    var found bool

    // Cek di AdminPerpustakaan
    var adminPerpus models.AdminPerpustakaan
    if err := config.DB.Where("username = ? OR email = ?", req.Username, req.Username).
        Preload("Perpustakaan").First(&adminPerpus).Error; err == nil {
        if err := bcrypt.CompareHashAndPassword([]byte(adminPerpus.Password), []byte(req.Password)); err == nil {
            if !adminPerpus.IsActive {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun belum diverifikasi"})
                return
            }
            
            adminPerpus.LastLogin = &now
            config.DB.Save(&adminPerpus)

            token, expiresAt, err := middleware.GenerateToken(adminPerpus.ID, adminPerpus.Username, "admin_perpustakaan")
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
                return
            }

            response = models.LoginResponse{
                Token:     token,
                User:      adminPerpus,
                UserType:  "admin_perpustakaan",
                ExpiresAt: expiresAt,
            }
            found = true
        }
    }

    // Cek di AdminDPK jika belum ditemukan
    if !found {
        var adminDPK models.AdminDPK
        if err := config.DB.Where("username = ? OR email = ?", req.Username, req.Username).
            First(&adminDPK).Error; err == nil {
            if err := bcrypt.CompareHashAndPassword([]byte(adminDPK.Password), []byte(req.Password)); err == nil {
                adminDPK.LastLogin = &now
                config.DB.Save(&adminDPK)

                token, expiresAt, err := middleware.GenerateToken(adminDPK.ID, adminDPK.Username, "admin_dpk")
                if err != nil {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
                    return
                }

                response = models.LoginResponse{
                    Token:     token,
                    User:      adminDPK,
                    UserType:  "admin_dpk",
                    ExpiresAt: expiresAt,
                }
                found = true
            }
        }
    }

    // Cek di Executive jika belum ditemukan
    if !found {
        var executive models.Executive
        if err := config.DB.Where("username = ? OR email = ?", req.Username, req.Username).
            First(&executive).Error; err == nil {
            if err := bcrypt.CompareHashAndPassword([]byte(executive.Password), []byte(req.Password)); err == nil {
                executive.LastLogin = &now
                config.DB.Save(&executive)

                token, expiresAt, err := middleware.GenerateToken(executive.ID, executive.Username, "executive")
                if err != nil {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
                    return
                }

                response = models.LoginResponse{
                    Token:     token,
                    User:      executive,
                    UserType:  "executive",
                    ExpiresAt: expiresAt,
                }
                found = true
            }
        }
    }

    if !found {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/email atau password salah"})
        return
    }

    // Create audit log
    var userID uint
    switch u := response.User.(type) {
    case models.AdminPerpustakaan:
        userID = u.ID
    case models.AdminDPK:
        userID = u.ID
    case models.Executive:
        userID = u.ID
    }

    auditLog := models.AuditLog{
        UserType:  response.UserType,
        UserID:    userID,
        Action:    "LOGIN",
        IPAddress: c.ClientIP(),
        UserAgent: c.GetHeader("User-Agent"),
        Timestamp: now,
    }
    config.DB.Create(&auditLog)

    c.JSON(http.StatusOK, response)
}

func Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	// Create audit log
	auditLog := models.AuditLog{
		UserType:  userType,
		UserID:    userID,
		Action:    "LOGOUT",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func ChangePassword(c *gin.Context) {
	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	var userModel interface{}
	var db *gorm.DB = config.DB

	switch userType {
	case "admin_perpustakaan":
		var adminPerpus models.AdminPerpustakaan
		if err := db.First(&adminPerpus, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(adminPerpus.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah"})
			return
		}
		
		adminPerpus.Password = string(hashedPassword)
		userModel = &adminPerpus

	case "admin_dpk":
		var adminDPK models.AdminDPK
		if err := db.First(&adminDPK, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(adminDPK.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah"})
			return
		}
		
		adminDPK.Password = string(hashedPassword)
		userModel = &adminDPK

	case "executive":
		var executive models.Executive
		if err := db.First(&executive, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(executive.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah"})
			return
		}
		
		executive.Password = string(hashedPassword)
		userModel = &executive

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe user tidak valid"})
		return
	}

	if err := db.Save(userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserType:  userType,
		UserID:    userID,
		Action:    "CHANGE_PASSWORD",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("UserAgent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

func GetProfile(c *gin.Context) {
    userID := c.GetUint("user_id")
    userType := c.GetString("user_type")

    var response gin.H

    switch userType {
    case "admin_perpustakaan":
        var adminPerpus models.AdminPerpustakaan
        if err := config.DB.Preload("Perpustakaan").First(&adminPerpus, userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
            return
        }
        response = gin.H{
            "id": adminPerpus.ID,
            "username": adminPerpus.Username,
            "nama_lengkap": adminPerpus.NamaLengkap,
            "email": adminPerpus.Email,
            "no_telepon": adminPerpus.NoTelepon,
			"user_type": "admin_perpustakaan",
            "perpustakaan": gin.H{
                "id": adminPerpus.Perpustakaan.ID,
                "nama_perpustakaan": adminPerpus.Perpustakaan.NamaPerpustakaan,
                "alamat": adminPerpus.Perpustakaan.Alamat,
                "jenis_perpustakaan": adminPerpus.Perpustakaan.JenisPerpustakaan,
            },
        }

    case "admin_dpk":
        var adminDPK models.AdminDPK
        if err := config.DB.First(&adminDPK, userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
            return
        }
        response = gin.H{
            "id": adminDPK.ID,
            "username": adminDPK.Username,
            "nama_lengkap": adminDPK.NamaLengkap,
            "email": adminDPK.Email,
            "no_telepon": adminDPK.NoTelepon,
        }

    case "executive":
        var executive models.Executive
        if err := config.DB.First(&executive, userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
            return
        }
        response = gin.H{
            "id": executive.ID,
            "username": executive.Username,
            "nama_lengkap": executive.NamaLengkap,
            "email": executive.Email,
            "no_telepon": executive.NoTelepon,
        }

    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe user tidak valid"})
        return
    }

    c.JSON(http.StatusOK, response)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove fields that shouldn't be updated
	delete(updateData, "id")
	delete(updateData, "username")
	delete(updateData, "password")
	delete(updateData, "is_active")
	delete(updateData, "last_login")
	delete(updateData, "created_at")
	delete(updateData, "updated_at")
	delete(updateData, "deleted_at")

	var err error
	var user interface{}

	switch userType {
	case "admin_perpustakaan":
		var adminPerpus models.AdminPerpustakaan
		if err = config.DB.First(&adminPerpus, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		
		if err = config.DB.Model(&adminPerpus).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		user = adminPerpus

	case "admin_dpk":
		var adminDPK models.AdminDPK
		if err = config.DB.First(&adminDPK, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		
		if err = config.DB.Model(&adminDPK).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		user = adminDPK

	case "executive":
		var executive models.Executive
		if err = config.DB.First(&executive, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		
		if err = config.DB.Model(&executive).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		user = executive

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe user tidak valid"})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserType:  userType,
		UserID:    userID,
		Action:    "UPDATE_PROFILE",
		TableName: userType,
		RecordID:  userID,
		NewValues: stringifyMap(updateData),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil diperbarui",
		"user":    user,
	})
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validasi gagal: " + err.Error()})
		return
	}

	// Only allow registration for admin_perpustakaan
	if req.UserType != "admin_perpustakaan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya admin perpustakaan yang dapat mendaftar melalui endpoint ini"})
		return
	}

	// Check for duplicates before starting transaction
	var existingAdminPerpus models.AdminPerpustakaan
	if err := config.DB.Where("username = ?", req.Username).First(&existingAdminPerpus).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	if err := config.DB.Where("email = ?", req.Email).First(&existingAdminPerpus).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	var existingPerpustakaan models.Perpustakaan
	if err := config.DB.Where("nomor_induk = ?", req.NomorInduk).First(&existingPerpustakaan).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor induk perpustakaan sudah digunakan"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password: " + err.Error()})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal"})
		}
	}()
	// Convert NomorInduk from string to int
	nomorIndukInt, err := strconv.Atoi(req.NomorInduk)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor induk harus berupa angka"})
		return
	}
	// Create new Perpustakaan
	perpustakaan := models.Perpustakaan{
		NamaPerpustakaan:  req.NamaPerpustakaan,
		Alamat:            req.Alamat,
		JenisPerpustakaan: req.JenisPerpustakaan,
		NomorInduk:        nomorIndukInt,
		StatusVerifikasi:  "Pending",
	}

	if err := tx.Create(&perpustakaan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat perpustakaan: " + err.Error()})
		return
	}

	// Create new AdminPerpustakaan
	admin := models.AdminPerpustakaan{
		Username:       req.Username,
		Password:       string(hashedPassword),
		NamaLengkap:    req.NamaLengkap,
		Email:          req.Email,
		NoTelepon:      req.NoTelepon,
		PerpustakaanID: perpustakaan.ID,
		IsActive:       false, // Inactive until verified
	}

	if err := tx.Create(&admin).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat admin perpustakaan: " + err.Error()})
		return
	}

	// Create verification entry
	verifikasi := models.Verifikasi{
		PerpustakaanID:    perpustakaan.ID,
		JenisData:         "Perpustakaan",
		Status:            "Pending",
		TanggalVerifikasi: nil,
		AdminDPKID: 0,
	}

	if err := tx.Create(&verifikasi).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat entri verifikasi: " + err.Error()})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_perpustakaan",
		UserID:    admin.ID,
		Action:    "REGISTER",
		TableName: "admin_perpustakaan",
		RecordID:  admin.ID,
		NewValues: stringifyMap(map[string]interface{}{
			
			"perpustakaan_id": perpustakaan.ID,
		}),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	if err := tx.Create(&auditLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat audit log: " + err.Error()})
		return
	}

	// Create notification for admin_dpk
	notifikasi := models.Notifikasi{
		JudulNotifikasi: "Registrasi Perpustakaan Baru",
		IsiNotifikasi:   "Perpustakaan baru dengan nama " + req.NamaPerpustakaan + " telah terdaftar dan menunggu verifikasi.",
		JenisNotifikasi: "info",
		TujuanUser:      "admin_dpk",
		RelatedID:       &perpustakaan.ID,
		RelatedType:     "perpustakaan",
		IsRead:          false,
		TanggalKirim:    time.Now(),
		ExpiredAt:       time.Now().Add(7 * 24 * time.Hour), // Expires in 7 days
		ActionLink:      "/admin-dpk/pending-verification",
	}
	if err := tx.Create(&notifikasi).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat notifikasi: " + err.Error()})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan transaksi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registrasi berhasil, menunggu verifikasi dari admin DPK",
		"user":    admin,
		"perpustakaan": perpustakaan,
	})
}

func RegisterAdminDPK(c *gin.Context) {
	var req models.RegisterAdminDPKRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only allow registration for admin_dpk
	if req.UserType != "admin_dpk" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya admin DPK yang dapat mendaftar melalui endpoint ini"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// Check if username already exists
	var existingAdminDPK models.AdminDPK
	if err := config.DB.Where("username = ?", req.Username).First(&existingAdminDPK).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	// Check if email already exists
	if err := config.DB.Where("email = ?", req.Email).First(&existingAdminDPK).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	// Create new AdminDPK
	admin := models.AdminDPK{
		Username:    req.Username,
		Password:    string(hashedPassword),
		NamaLengkap: req.NamaLengkap,
		Email:       req.Email,
		NoTelepon:   req.NoTelepon,
		HakAkses:    "admin",
		IsActive:    true, // Active immediately
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat admin DPK"})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    admin.ID,
		Action:    "REGISTER",
		TableName: "admin_dpk",
		RecordID:  admin.ID,
		NewValues: stringifyMap(map[string]interface{}{
			"username":     req.Username,
			"nama_lengkap": req.NamaLengkap,
			"email":        req.Email,
		}),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{
		"message": "Registrasi admin DPK berhasil",
		"user":    admin,
	})
}

func stringifyMap(m map[string]interface{}) string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(bytes)
}
func RegisterExecutive(c *gin.Context) {
    var req models.RegisterAdminDPKRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if req.UserType != "executive" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya executive yang dapat mendaftar melalui endpoint ini"})
        return
    }
    var existingExecutive models.Executive
    if err := config.DB.Where("username = ?", req.Username).First(&existingExecutive).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
        return
    }
    if err := config.DB.Where("email = ?", req.Email).First(&existingExecutive).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
        return
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
        return
    }
    executive := models.Executive{
        Username:    req.Username,
        Password:    string(hashedPassword),
        NamaLengkap: req.NamaLengkap,
        Email:       req.Email,
        NoTelepon:   req.NoTelepon,
        Jabatan:     "Executive",
        IsActive:    true,
    }
    if err := config.DB.Create(&executive).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat executive"})
        return
    }
    auditLog := models.AuditLog{
        UserType:  "executive",
        UserID:    executive.ID,
        Action:    "REGISTER",
        TableName: "executive",
        RecordID:  executive.ID,
        NewValues: stringifyMap(map[string]interface{}{
            "username":     req.Username,
            "nama_lengkap": req.NamaLengkap,
            "email":        req.Email,
        }),
        IPAddress: c.ClientIP(),
        UserAgent: c.GetHeader("User-Agent"),
        Timestamp: time.Now(),
    }
    config.DB.Create(&auditLog)
    c.JSON(http.StatusOK, gin.H{
        "message": "Registrasi executive berhasil",
        "user":    executive,
    })
}

func GetPendingAdminVerifications(c *gin.Context) {
    var pendingAdmins []models.AdminPerpustakaan
    
    // Get all admin perpustakaan with is_active = false
    if err := config.DB.Preload("Perpustakaan").
        Where("is_active = ?", false).
        Find(&pendingAdmins).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data admin pending"})
        return
    }
    
    c.JSON(http.StatusOK, pendingAdmins)
}

func VerifyAdminPerpustakaan(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    var req models.VerifyAdminRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Start transaction
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal"})
        }
    }()

    // Get admin perpustakaan data
    var adminPerpus models.AdminPerpustakaan
    if err := tx.Preload("Perpustakaan").First(&adminPerpus, req.AdminPerpustakaanID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
        return
    }

    // Update admin status
    isApproved := req.Status == "approved"
    adminPerpus.IsActive = isApproved
    if err := tx.Save(&adminPerpus).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status admin"})
        return
    }

    // Update perpustakaan verification status
    if isApproved {
        adminPerpus.Perpustakaan.StatusVerifikasi = "Verified"
        if err := tx.Save(&adminPerpus.Perpustakaan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status perpustakaan"})
            return
        }
    }

    // Create verification record
    verifikasi := models.Verifikasi{
        PerpustakaanID:    adminPerpus.PerpustakaanID,
        JenisData:        "AdminPerpustakaan",
        Status:           req.Status,
        CatatanRevisi:    req.Catatan,
		TanggalVerifikasi: func() *time.Time { t := time.Now(); return &t }(),
        AdminDPKID:      userID,
    }
    if err := tx.Create(&verifikasi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat record verifikasi"})
        return
    }

    // Create notification for admin perpustakaan
    judul := "Akun Anda Telah Diverifikasi"
    isi := "Akun admin perpustakaan Anda telah diverifikasi dan diaktifkan."
    if !isApproved {
        judul = "Verifikasi Akun Ditolak"
        isi = "Verifikasi akun admin perpustakaan Anda ditolak. Catatan: " + req.Catatan
    }

    notifikasi := models.Notifikasi{
        JudulNotifikasi: judul,
        IsiNotifikasi:   isi,
        JenisNotifikasi: "info",
        TujuanUser:      "admin_perpustakaan",
        UserID:         &adminPerpus.ID,
        RelatedID:      &adminPerpus.ID,
        RelatedType:    "admin_perpustakaan",
        IsRead:         false,
        TanggalKirim:   time.Now(),
        ExpiredAt:      time.Now().Add(7 * 24 * time.Hour),
        ActionLink:     "/admin-perpustakaan/profile",
    }
    if err := tx.Create(&notifikasi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat notifikasi"})
        return
    }

    // Create audit log
    auditLog := models.AuditLog{
        UserType:  "admin_dpk",
        UserID:    userID,
        Action:    "VERIFY_ADMIN_PERUSTAKAAN",
        TableName: "admin_perpustakaan",
        RecordID:  adminPerpus.ID,
        NewValues: stringifyMap(map[string]interface{}{
            "status":        req.Status,
            "is_active":     isApproved,
            "admin_dpk_id":  userID,
            "catatan":       req.Catatan,
        }),
        IPAddress: c.ClientIP(),
        UserAgent: c.GetHeader("User-Agent"),
        Timestamp: time.Now(),
    }
    if err := tx.Create(&auditLog).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat audit log"})
        return
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan transaksi"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Verifikasi admin perpustakaan berhasil",
        "status":  req.Status,
        "admin":   adminPerpus,
    })
}