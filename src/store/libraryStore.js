  import { defineStore } from 'pinia'
  import { ref } from 'vue'
  import api from '../api/axios'

  export const useLibraryStore = defineStore('library', () => {
    const libraries = ref([])
    const currentLibrary = ref(null)

    // Fungsi untuk menambah data perpustakaan baru ke backend
    const addLibrary = async (libraryData) => {
      const payload = {
        periode: libraryData.periode,
        nama_perpustakaan: libraryData.nama_perpustakaan,
        alamat: libraryData.alamat,
        kepala_perpustakaan: libraryData.kepala_perpustakaan,
        jenis_perpustakaan: libraryData.jenis_perpustakaan,
        tahun_berdiri: libraryData.tahun_berdiri,
        nomor_induk: libraryData.nomor_induk,
        jumlah_sdm: libraryData.jumlah_sdm,
        jumlah_pengunjung: libraryData.jumlah_pengunjung,
        jumlah_anggota: libraryData.jumlah_anggota
      }

      const response = await api.post('/admin-perpustakaan/input-data', payload)
      await fetchLibraries()
      return response.data
    }

    // Fungsi untuk mendapatkan data perpustakaan berdasarkan ID
    const getLibraryById = async (id) => {
      try {
        const response = await api.get(`/admin-perpustakaan/data/${id}`)
        return response.data
      } catch (error) {
        console.error('Error getting library:', error)
        throw error
      }
    }

    // Ambil data perpustakaan dari backend (semua milik user login)
    const fetchLibraries = async () => {
    try {
      const response = await api.get('/admin-perpustakaan/data')
      
      // Handle berbagai format response
      let data = []
      if (Array.isArray(response.data)) {
        data = response.data
      } else if (response.data && typeof response.data === 'object') {
        if (Array.isArray(response.data.data)) {
          data = response.data.data
        } else if (Array.isArray(response.data.result)) {
          data = response.data.result
        } else if (response.data.id) { // Jika single object
          data = [response.data]
        }
      }

      libraries.value = data.map(item => ({
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
    } catch (error) {
      console.error('Error fetching libraries:', error)
      throw error
    }
  }

    // Hapus data perpustakaan dari backend
    const deleteLibrary = async (id) => {
      try {
        await api.delete(`/admin-perpustakaan/data/${id}`)
        await fetchLibraries() // Refresh data setelah menghapus
        return true
      } catch (error) {
        console.error('Error deleting library:', error)
        throw error
      }
    }

    // Update data perpustakaan ke backend
    const updateLibrary = async (id, updatedData) => {
      try {
        const response = await api.put(`/admin-perpustakaan/data/${id}`, updatedData)
        await fetchLibraries() // Refresh data setelah mengupdate
        return response.data
      } catch (error) {
        console.error('Error updating library:', error)
        throw error
      }
    }
    const sendDataToDPK = async (id) => {
      const payload = {
        perpustakaan_id: id,
        catatan_kirim: '' // Bisa dikosongkan atau nanti dikasih catatan manual
      }

      const response = await api.post(`/admin-perpustakaan/data/${id}/send-data`, payload)
      await fetchLibraries()
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