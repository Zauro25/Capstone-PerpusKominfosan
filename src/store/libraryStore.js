import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLibraryStore = defineStore('library', () => {
  const libraries = ref([])
  const currentLibrary = ref(null)

  // Fungsi untuk menambah data perpustakaan baru
  const addLibrary = (libraryData) => {
    const newLibrary = {
      id: Date.now(), // Temporary ID generator
      ...libraryData
    }
    libraries.value.push(newLibrary)
    return newLibrary
  }

  // Fungsi untuk mendapatkan data perpustakaan berdasarkan ID
  const getLibraryById = (id) => {
    return libraries.value.find(lib => lib.id === id)
  }

  // Fungsi untuk memperbarui data perpustakaan
  const updateLibrary = (id, updatedData) => {
    const index = libraries.value.findIndex(lib => lib.id === id)
    if (index !== -1) {
      libraries.value[index] = {
        ...libraries.value[index],
        ...updatedData
      }
      return libraries.value[index]
    }
    return null
  }

  // Fungsi untuk menghapus data perpustakaan
  const deleteLibrary = (id) => {
    const index = libraries.value.findIndex(lib => lib.id === id)
    if (index !== -1) {
      libraries.value.splice(index, 1)
      return true
    }
    return false
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
    setCurrentLibrary
  }
}) 