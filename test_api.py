import requests
import json
from fpdf import FPDF
from datetime import datetime
import time

# Konfigurasi
BASE_URL = "http://localhost:8080/api/v1"
REPORT_FILE = "api_test_report.pdf"
HEADERS = {'Content-Type': 'application/json'}

# Inisialisasi PDF
pdf = FPDF()
pdf.add_page()
pdf.set_font("Arial", size=12)

# Fungsi untuk menambahkan hasil test ke PDF
def add_to_report(title, description, response, success):
    pdf.set_font("Arial", 'B', 14)
    pdf.cell(200, 10, txt=title, ln=1)
    pdf.set_font("Arial", size=12)
    
    pdf.cell(200, 10, txt=f"Endpoint: {description}", ln=1)
    pdf.cell(200, 10, txt=f"Status Code: {response.status_code}", ln=1)
    
    if success:
        pdf.set_text_color(0, 128, 0)  # Hijau untuk sukses
        pdf.cell(200, 10, txt="RESULT: SUCCESS", ln=1)
    else:
        pdf.set_text_color(255, 0, 0)  # Merah untuk gagal
        pdf.cell(200, 10, txt="RESULT: FAILED", ln=1)
    
    try:
        response_data = response.json()
        pdf.multi_cell(0, 10, txt=f"Response: {json.dumps(response_data, indent=2)}")
    except:
        pdf.multi_cell(0, 10, txt=f"Response: {response.text}")
    
    pdf.set_text_color(0, 0, 0)  # Kembali ke hitam
    pdf.ln(5)

# Fungsi untuk membersihkan database
def clear_database():
    # Hanya contoh - dalam implementasi nyata, ini akan memanggil endpoint khusus atau script database
    print("Membersihkan database...")
    # Di sini Anda bisa menambahkan kode untuk menghapus data testing
    
    # Contoh: Hapus semua admin perpustakaan
    try:
        login_response = requests.post(
            f"{BASE_URL}/login",
            json={"username": "admin_dpk", "password": "admin123", "user_type": "admin_dpk"},
            headers=HEADERS
        )
        token = login_response.json().get('token')
        
        headers = HEADERS.copy()
        headers['Authorization'] = f"Bearer {token}"
        
        # Hapus semua admin perpustakaan
        response = requests.get(f"{BASE_URL}/admin-dpk/manage-accounts", headers=headers)
        admins = response.json()
        for admin in admins:
            requests.delete(f"{BASE_URL}/admin-dpk/manage-accounts/{admin['id']}", headers=headers)
        
        add_to_report("Clear Database", "Menghapus data testing", response, True)
    except Exception as e:
        add_to_report("Clear Database", f"Gagal menghapus data testing: {str(e)}", 
                     type('obj', (object,), {'status_code': 500, 'text': str(e)}), False)

# Fungsi untuk test register dan login semua tipe user
def test_user_flows():
    test_results = []
    
    # Data testing
    test_users = [
        {
            "type": "admin_perpustakaan",
            "register_data": {
                "user_type": "admin_perpustakaan",
                "username": "test_admin_perpus",
                "password": "test123",
                "nama_lengkap": "Test Admin Perpus",
                "email": "test_admin_perpus@example.com",
                "no_telepon": "081234567890",
                "nama_perpustakaan": "Perpustakaan Test",
                "alamat": "Jl. Test No. 123",
                "jenis_perpustakaan": "Umum",
                "nomor_induk": "TEST123"
            },
            "login_data": {
                "username": "test_admin_perpus",
                "password": "test123",
                "user_type": "admin_perpustakaan"
            }
        },
        {
            "type": "admin_dpk",
            "register_data": {
                "user_type": "admin_dpk",
                "username": "test_admin_dpk",
                "password": "test123",
                "nama_lengkap": "Test Admin DPK",
                "email": "test_admin_dpk@example.com",
                "no_telepon": "081234567891"
            },
            "login_data": {
                "username": "test_admin_dpk",
                "password": "test123",
                "user_type": "admin_dpk"
            }
        },
        {
            "type": "executive",
            "register_data": {
                "user_type": "executive",
                "username": "test_executive",
                "password": "test123",
                "nama_lengkap": "Test Executive",
                "email": "test_executive@example.com",
                "no_telepon": "081234567892"
            },
            "login_data": {
                "username": "test_executive",
                "password": "test123",
                "user_type": "executive"
            }
        }
    ]
    
    # Test register dan login untuk setiap tipe user
    for user in test_users:
        # Test register
        try:
            response = requests.post(
                f"{BASE_URL}/register",
                json=user["register_data"],
                headers=HEADERS
            )
            success = response.status_code in [200, 201]
            test_results.append((
                f"Register {user['type']}",
                f"POST {BASE_URL}/register",
                response,
                success
            ))
            
            # Jika admin perpustakaan, perlu diverifikasi dulu
            if user["type"] == "admin_perpustakaan" and success:
                # Login sebagai admin DPK untuk verifikasi
                login_response = requests.post(
                    f"{BASE_URL}/login",
                    json={"username": "admin_dpk", "password": "admin123", "user_type": "admin_dpk"},
                    headers=HEADERS
                )
                token = login_response.json().get('token')
                
                headers = HEADERS.copy()
                headers['Authorization'] = f"Bearer {token}"
                
                # Dapatkan ID admin perpustakaan yang baru dibuat
                admin_id = response.json().get('user', {}).get('id')
                
                # Verifikasi admin perpustakaan
                verify_response = requests.post(
                    f"{BASE_URL}/admin-dpk/verify-admin-perpustakaan",
                    json={"admin_perpustakaan_id": admin_id, "status": "approved"},
                    headers=headers
                )
                test_results.append((
                    f"Verify {user['type']}",
                    f"POST {BASE_URL}/admin-dpk/verify-admin-perpustakaan",
                    verify_response,
                    verify_response.status_code in [200, 201]
                ))
        except Exception as e:
            test_results.append((
                f"Register {user['type']}",
                f"POST {BASE_URL}/register",
                type('obj', (object,), {'status_code': 500, 'text': str(e)}),
                False
            ))
        
        # Test login
        try:
            response = requests.post(
                f"{BASE_URL}/login",
                json=user["login_data"],
                headers=HEADERS
            )
            test_results.append((
                f"Login {user['type']}",
                f"POST {BASE_URL}/login",
                response,
                response.status_code == 200
            ))
        except Exception as e:
            test_results.append((
                f"Login {user['type']}",
                f"POST {BASE_URL}/login",
                type('obj', (object,), {'status_code': 500, 'text': str(e)}),
                False
            ))
    
    # Tambahkan hasil ke PDF
    for title, desc, response, success in test_results:
        add_to_report(title, desc, response, success)

# Fungsi untuk test endpoint admin perpustakaan
def test_admin_perpustakaan_endpoints():
    # Login sebagai admin perpustakaan
    login_response = requests.post(
        f"{BASE_URL}/login",
        json={"username": "test_admin_perpus", "password": "test123", "user_type": "admin_perpustakaan"},
        headers=HEADERS
    )
    token = login_response.json().get('token')
    
    headers = HEADERS.copy()
    headers['Authorization'] = f"Bearer {token}"
    
    # Test endpoint
    endpoints = [
        ("GET Profile", f"{BASE_URL}/profile", "GET", None),
        ("Update Profile", f"{BASE_URL}/profile", "PUT", {"nama_lengkap": "Updated Name"}),
        ("Change Password", f"{BASE_URL}/change-password", "POST", 
         {"old_password": "test123", "new_password": "newpassword123"}),
        ("Get Dashboard", f"{BASE_URL}/admin-perpustakaan/dashboard", "GET", None),
        ("Get Data Perpustakaan", f"{BASE_URL}/admin-perpustakaan/data", "GET", None),
        ("Update Data Perpustakaan", f"{BASE_URL}/admin-perpustakaan/data", "PUT", 
         {"alamat": "Jl. Updated No. 456"}),
        ("Send Data to DPK", f"{BASE_URL}/admin-perpustakaan/send-data", "POST", 
         {"perpustakaan_id": 1, "catatan_kirim": "Data untuk verifikasi"}),
        ("Get History", f"{BASE_URL}/admin-perpustakaan/history", "GET", None),
        ("Get Notifications", f"{BASE_URL}/notifications", "GET", None)
    ]
    
    for title, url, method, data in endpoints:
        try:
            if method == "GET":
                response = requests.get(url, headers=headers)
            elif method == "POST":
                response = requests.post(url, json=data, headers=headers)
            elif method == "PUT":
                response = requests.put(url, json=data, headers=headers)
            
            add_to_report(title, f"{method} {url}", response, response.status_code in [200, 201])
        except Exception as e:
            add_to_report(title, f"{method} {url}", 
                         type('obj', (object,), {'status_code': 500, 'text': str(e)}), False)

# Fungsi untuk test endpoint admin DPK
def test_admin_dpk_endpoints():
    # Login sebagai admin DPK
    login_response = requests.post(
        f"{BASE_URL}/login",
        json={"username": "test_admin_dpk", "password": "test123", "user_type": "admin_dpk"},
        headers=HEADERS
    )
    token = login_response.json().get('token')
    
    headers = HEADERS.copy()
    headers['Authorization'] = f"Bearer {token}"
    
    # Test endpoint
    endpoints = [
        ("GET All Perpustakaan", f"{BASE_URL}/admin-dpk/perpustakaan", "GET", None),
        ("GET Detail Perpustakaan", f"{BASE_URL}/admin-dpk/perpustakaan/1", "GET", None),
        ("Verify Data", f"{BASE_URL}/admin-dpk/verifikasi", "POST", 
         {"perpustakaan_id": 1, "status": "approved", "catatan_revisi": "Data sudah valid"}),
        ("Get Pending Verification", f"{BASE_URL}/admin-dpk/pending-verification", "GET", None),
        ("Generate Report", f"{BASE_URL}/admin-dpk/laporan", "POST", 
         {"judul": "Laporan Test", "periode": "2023-1", "jenis_laporan": "Statistik", "format_file": "PDF"}),
        ("Get All Reports", f"{BASE_URL}/admin-dpk/laporan", "GET", None),
        ("Broadcast Notification", f"{BASE_URL}/admin-dpk/notifications/broadcast", "POST", 
         {"judul_notifikasi": "Pengingat", "isi_notifikasi": "Jangan lupa update data", 
          "jenis_notifikasi": "warning", "tujuan_user": "admin_perpustakaan"}),
        ("Get Audit Logs", f"{BASE_URL}/admin-dpk/audit-logs", "GET", None),
        ("Manage Admin Accounts", f"{BASE_URL}/admin-dpk/manage-accounts", "GET", None)
    ]
    
    for title, url, method, data in endpoints:
        try:
            if method == "GET":
                response = requests.get(url, headers=headers)
            elif method == "POST":
                response = requests.post(url, json=data, headers=headers)
            
            add_to_report(title, f"{method} {url}", response, response.status_code in [200, 201])
        except Exception as e:
            add_to_report(title, f"{method} {url}", 
                         type('obj', (object,), {'status_code': 500, 'text': str(e)}), False)

# Fungsi untuk test endpoint executive
def test_executive_endpoints():
    # Login sebagai executive
    login_response = requests.post(
        f"{BASE_URL}/login",
        json={"username": "test_executive", "password": "test123", "user_type": "executive"},
        headers=HEADERS
    )
    token = login_response.json().get('token')
    
    headers = HEADERS.copy()
    headers['Authorization'] = f"Bearer {token}"
    
    # Test endpoint
    endpoints = [
        ("GET Dashboard", f"{BASE_URL}/executive/dashboard", "GET", None),
        ("GET Statistics", f"{BASE_URL}/executive/statistics", "GET", None),
        ("GET Reports", f"{BASE_URL}/executive/laporan", "GET", None)
    ]
    
    for title, url, method, data in endpoints:
        try:
            if method == "GET":
                response = requests.get(url, headers=headers)
            
            add_to_report(title, f"{method} {url}", response, response.status_code in [200, 201])
        except Exception as e:
            add_to_report(title, f"{method} {url}", 
                         type('obj', (object,), {'status_code': 500, 'text': str(e)}), False)

# Fungsi utama
def main():
    # Header PDF
    pdf.set_font("Arial", 'B', 16)
    pdf.cell(200, 10, txt="API TEST REPORT", ln=1, align='C')
    pdf.set_font("Arial", size=12)
    pdf.cell(200, 10, txt=f"Tanggal: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}", ln=1)
    pdf.ln(10)
    
    print("Memulai testing API...")
    
    # 1. Bersihkan database
    clear_database()
    
    # 2. Test alur user (register, login)
    print("Testing alur user...")
    test_user_flows()
    
    # 3. Test endpoint admin perpustakaan
    print("Testing endpoint admin perpustakaan...")
    test_admin_perpustakaan_endpoints()
    
    # 4. Test endpoint admin DPK
    print("Testing endpoint admin DPK...")
    test_admin_dpk_endpoints()
    
    # 5. Test endpoint executive
    print("Testing endpoint executive...")
    test_executive_endpoints()
    
    # Simpan PDF
    pdf.output(REPORT_FILE)
    print(f"Testing selesai. Laporan disimpan di {REPORT_FILE}")

if __name__ == "__main__":
    main()