import { AxiosResponse } from "axios";
import { DIRECT } from "../constant/constant";
import { apiClient } from "./axios";
import {
  CheckSessionDTO,
  ValidateDTO,
  VerifyGoogleLoginDTO,
} from "./dto/auth.dto";
import {
  parseCheckSession,
  parseCredentials,
  parseVerifyGoogleLogin,
} from "./parser/auth.parser";

export const checkSession = async (serviceUrl: string | null) => {
  try {
    const res: AxiosResponse<CheckSessionDTO> = await apiClient.get(
      "/auth/check-session",
      {
        params: { service: serviceUrl ?? DIRECT },
        withCredentials: true,
      }
    );

    return parseCheckSession(res.data);
  } catch {
    const defaultResp = { serviceTicket: "" };
    return defaultResp;
  }
};

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
        withCredentials: true,
      }
    );

    return parseVerifyGoogleLogin(res.data);
  } catch (error) {
    console.error("Failed to verify Google login: ", error);

    return Error("Failed to verify Google login");
  }
};

export const validateST = async (serviceTicket: string) => {
  try {
    const res: AxiosResponse<ValidateDTO> = await apiClient.get(
      "/auth/validate-st",
      {
        params: {
          ticket: serviceTicket,
          service: DIRECT,
        },
      }
    );

    return parseCredentials(res.data);
  } catch {
    return Error("Failed to validate service ticket");
  }
};

export const signout = async () => {
  try {
    const res: AxiosResponse<null> = await apiClient.post(
      "/auth/signout",
      {},
      { withCredentials: true }
    );

    return res.data;
  } catch {
    return Error("Failed to signout");
  }
};
