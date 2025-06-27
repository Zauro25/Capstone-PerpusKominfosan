import requests
import json
import psycopg2
import time
import os
BASE_URL = "http://localhost:8080"

# Store tokens and IDs for reuse
tokens = {
    "admin_perpustakaan": None,
    "admin_dpk": None,
    "executive": None
}
ids = {
    "admin_perpustakaan": None,
    "perpustakaan": None,
    "admin_dpk": None,
    "executive": None,
    "notification": None,
    "laporan": None,
    "admin_perpustakaan_2": None
}

# Generate unique suffix based on timestamp
unique_suffix = str(int(time.time()))

# File to store test results
LOG_FILE = "hasil.txt"
def clean_database():
    try:
        conn = psycopg2.connect(
            dbname="perpustakaan_db",
            user="perpustakaan_user",
            password="password123",
            host="localhost",
            port="5432"
        )
        cur = conn.cursor()
        # List tabel yang mungkin ada
        tables = [
            "audit_logs", "notifikasis", "log_revisis", "laporans", "verifikasis",
            "pengunjungs", "sdms", "koleksis", "admin_perpustakaans",
            "perpustakaans", "admin_dpks", "executives", "anggota"  # Ganti anggotas dengan anggota
        ]
        # Cek tabel yang ada
        cur.execute("""
            SELECT table_name 
            FROM information_schema.tables 
            WHERE table_schema = 'public';
        """)
        existing_tables = [row[0] for row in cur.fetchall()]
        # Hapus data hanya dari tabel yang ada
        for table in tables:
            if table in existing_tables:
                cur.execute(f"DELETE FROM {table};")
                with open(LOG_FILE, "a") as f:
                    f.write(f"Cleared table {table}\n")
            else:
                with open(LOG_FILE, "a") as f:
                    f.write(f"Table {table} does not exist, skipping\n")
        conn.commit()
        with open(LOG_FILE, "a") as f:
            f.write("Database cleaned successfully!\n")
    except Exception as e:
        with open(LOG_FILE, "a") as f:
            f.write(f"Failed to clean database: {str(e)}\n")
    finally:
        cur.close()
        conn.close()


    print("\nAPI tests completed!")
    with open(LOG_FILE, "a") as f:
        f.write("\nAPI tests completed!\n")
    generate_latex()

if __name__ == "__main__":
    clean_database()
    run_tests()