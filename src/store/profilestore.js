import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/axios'

export const useProfileStore = defineStore('profile', () => {
  const userProfile = ref({
    id: null,
    username: '',
    nama_lengkap: '',
    user_type: '',
    email: '',
    no_telepon: '',
    perpustakaan: {
      id: null,
      nama_perpustakaan: '',
      alamat: '',
      jenis_perpustakaan: ''
    }
  })
  
  const loading = ref(false)
  const error = ref(null)

  const fetchUserProfile = async () => {
  loading.value = true;
  try {
    const response = await api.get('/profile');

    userProfile.value = {
      ...response.data,
      perpustakaan: response.data.perpustakaan || {} // Pastikan perpustakaan tidak null
    };
    console.log('Data dari API:', userProfile.value); // Debug
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

  const updateUserProfile = async (profileData) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.put('/profile', {
        nama_lengkap: profileData.nama_lengkap,
        email: profileData.email,
        no_telepon: profileData.no_telepon
      })
      userProfile.value = {
        ...userProfile.value,
        nama_lengkap: response.data.nama_lengkap,
        email: response.data.email,
        no_telepon: response.data.no_telepon
      }
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return { 
    userProfile, 
    loading, 
    error,
    fetchUserProfile, 
    updateUserProfile 
  }
})