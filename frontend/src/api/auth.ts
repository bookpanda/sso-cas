import { AxiosResponse } from "axios";
import { DIRECT } from "../constant/constant";
import { apiClient } from "./axios";
import { VerifyGoogleLoginDTO } from "./dto/auth.dto";
import { parseVerifyGoogleLogin } from "./parser/auth.parser";

export const getGoogleLoginUrl = async (serviceUrl: string | null) => {
  try {
    const res: AxiosResponse<string> = await apiClient.get("/auth/google-url", {
      params: { service: serviceUrl ?? DIRECT },
    });

    return res.data;
  } catch (error) {
    console.error("Failed to get Google login URL: ", error);

    return Error("Failed to get Google login URL");
  }
};

export const verifyGoogleLogin = async (code: string, state: string) => {
  try {
    const res: AxiosResponse<VerifyGoogleLoginDTO> = await apiClient.get(
      "/auth/verify-google",
      {
        params: { code: code, state: state },
      }
    );

    return parseVerifyGoogleLogin(res.data);
  } catch (error) {
    console.error("Failed to get Google login URL: ", error);

    return Error("Failed to get Google login URL");
  }
};
