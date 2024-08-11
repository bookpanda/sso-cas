import { AxiosResponse } from "axios";
import { apiClient } from "./axios";

export const getGoogleLoginUrl = async (serviceUrl: string) => {
  try {
    const res: AxiosResponse<string> = await apiClient.get("/auth/google-url", {
      params: { service: serviceUrl },
    });

    return res;
  } catch (error) {
    console.error("Failed to get Google login URL: ", error);

    return Error("Failed to get Google login URL");
  }
};
