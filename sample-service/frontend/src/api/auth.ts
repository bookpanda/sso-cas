import { AxiosResponse } from "axios";
import { DIRECT } from "../constant/constant";
import { apiClient } from "./axios";

export const getGoogleLoginUrl = async (serviceUrl: string | null) => {
  try {
    const res: AxiosResponse<string> = await apiClient.get("/auth/google-url", {
      params: { service: serviceUrl ?? DIRECT },
    });

    return res;
  } catch (error) {
    console.error("Failed to get Google login URL: ", error);

    return Error("Failed to get Google login URL");
  }
};

export const verifyGoogleLogin = async (code: string, state: string) => {
  try {
    const res: AxiosResponse<string> = await apiClient.get(
      "/auth/verify-google",
      {
        params: { code: code, state: state },
      }
    );

    return res;
  } catch (error) {
    console.error("Failed to get Google login URL: ", error);

    return Error("Failed to get Google login URL");
  }
};
