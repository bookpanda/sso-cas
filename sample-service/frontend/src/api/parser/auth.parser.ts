import { AuthenticateSSODTO } from "../dto/auth.dto";

export const parseAuthenticateSSO = (dto: AuthenticateSSODTO) => {
  return {
    accessToken: dto.access_token,
    refreshToken: dto.refresh_token,
    expiresIn: dto.expires_in,
  };
};
