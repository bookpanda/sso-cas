import axios from "axios";
import { API_URL } from "../constant/env";

export const apiClient = axios.create({
  baseURL: `${API_URL}/api/v1`,
  timeout: 10000,
});
