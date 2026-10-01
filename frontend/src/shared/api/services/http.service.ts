import axios from "axios";
import config from "../../configs/config.json"
import { STORAGE_KEYS } from "../../consts";

export const httpService = axios.create({
    baseURL: config.api.baseUrl
})

httpService.interceptors.request.use((config) => {
    const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
})