import { AxiosResponse } from "axios";
import { apiClient } from "./axios";
import { AuthenticateSSODTO } from "./dto/auth.dto";
import { parseAuthenticateSSO } from "./parser/auth.parser";

export const authenticateSSO = async (serviceTicket: string) => {
  try {
    const res: AxiosResponse<AuthenticateSSODTO> = await apiClient.get(
      "/auth/auth-sso",
      {
        params: { ticket: serviceTicket },
      }
    );

    return parseAuthenticateSSO(res.data);
  } catch (error) {
    console.error("Failed to authenticate SSO: ", error);

    return Error("Failed to authenticate SSO");
  }
};
