// utils/axiosConfig.ts

import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

const baseConfig: AxiosRequestConfig = {
  baseURL: process.env.API_BASE_URL || '',
  timeout: 5000, // Set your preferred timeout
};

export const publicClient: AxiosInstance = axios.create(baseConfig);

export const privateClient: AxiosInstance = axios.create(baseConfig);

// Interceptor for JWT-based requests
privateClient.interceptors.request.use((config) => {
  // Add JWT token to headers if available
  const token = localStorage.getItem('jwtToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default {
  publicClient,
  privateClient,
};
