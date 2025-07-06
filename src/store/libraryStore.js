import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/axios'

export const useLibraryStore = defineStore('library', () => {
  const libraries = ref([])
  const currentLibrary = ref(null)

  // Fungsi untuk menambah data perpustakaan baru ke backend
  const addLibrary = async (libraryData) => {
    // Mapping data dari form ke backend
    const payload = {
      periode: libraryData.periode,
      nama_perpustakaan: libraryData.namaPerpustakaan,
      alamat: libraryData.alamat,
      kepala_perpustakaan: libraryData.kepalaPerpustakaan,
      jenis_perpustakaan: libraryData.jenisPerpustakaan,
      tahun_berdiri: libraryData.tahunBerdiri,
      nomor_induk: libraryData.nomorInduk,
      jumlah_sdm: libraryData.jumlahSDM,
      jumlah_pengunjung: libraryData.jumlahPengunjung,
      jumlah_anggota: libraryData.jumlahAnggota
    }
    const response = await api.post('/admin-perpustakaan/input-data', payload)
    return response.data
  }

  // Fungsi untuk mendapatkan data perpustakaan berdasarkan ID
  const getLibraryById = (id) => {
    return libraries.value.find(lib => lib.id === id)
  }

  // Ambil data perpustakaan dari backend (semua milik user login)
  const fetchLibraries = async () => {
    const response = await api.get('/admin-perpustakaan/data-list')
    console.log('[fetchLibraries] response.data:', response.data)
    let dataArr = []
    if (Array.isArray(response.data)) {
      dataArr = response.data
    } else if (response.data && typeof response.data === 'object') {
      if ('id' in response.data) {
        dataArr = [response.data]
      } else if (Array.isArray(response.data.data)) {
        dataArr = response.data.data
      } else if (Array.isArray(response.data.result)) {
        dataArr = response.data.result
      } else {
        console.log('fetchLibraries fallback, response.data:', response.data)
      }
    } else {
      console.log('fetchLibraries fallback, response:', response)
    }
    console.log('dataArr:', dataArr)
    libraries.value = dataArr.map(item => ({
      id: item.id,
      periode: item.periode,
      nama_perpustakaan: item.nama_perpustakaan,
      alamat: item.alamat,
      kepala_perpustakaan: item.kepala_perpustakaan,
      jenis_perpustakaan: item.jenis_perpustakaan,
      tahun_berdiri: item.tahun_berdiri,
      nomor_induk: item.nomor_induk,
      jumlah_sdm: item.jumlah_sdm,
      jumlah_pengunjung: item.jumlah_pengunjung,
      jumlah_anggota: item.jumlah_anggota,
      status_verifikasi: item.status_verifikasi,
      tanggal_kirim: item.tanggal_kirim,
      catatan_revisi: item.catatan_revisi
    }))
    return libraries.value
  }

  // Hapus data perpustakaan dari backend
  const deleteLibrary = async (id) => {
    await api.delete(`/admin-perpustakaan/data/${id}`)
    // Optional: fetch ulang data setelah hapus
    return true
  }

  // Update data perpustakaan ke backend
  const updateLibrary = async (id, updatedData) => {
    // Pastikan mapping field ke snake_case
    const payload = {
      periode: updatedData.periode,
      nama_perpustakaan: updatedData.nama_perpustakaan,
      alamat: updatedData.alamat,
      kepala_perpustakaan: updatedData.kepala_perpustakaan,
      jenis_perpustakaan: updatedData.jenis_perpustakaan,
      tahun_berdiri: updatedData.tahun_berdiri,
      nomor_induk: updatedData.nomor_induk,
      jumlah_sdm: updatedData.jumlah_sdm,
      jumlah_pengunjung: updatedData.jumlah_pengunjung,
      jumlah_anggota: updatedData.jumlah_anggota
    }
    const response = await api.put(`/admin-perpustakaan/data/${id}`, payload)
    return response.data
  }

  // Fungsi untuk mengatur data perpustakaan yang sedang aktif
  const setCurrentLibrary = (library) => {
    currentLibrary.value = library
  }

  return {
    libraries,
    currentLibrary,
    addLibrary,
    getLibraryById,
    updateLibrary,
    deleteLibrary,
    setCurrentLibrary,
    fetchLibraries
  }
})