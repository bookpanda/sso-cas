export type AuthenticateSSODTO = {
  access_token: string;
  refresh_token: string;
  expires_in: Date;
};

export interface AuthToken {
  accessToken: string;
  refreshToken: string;
  expiresIn: Date;
}

export type ValidateDTO = {
  user_id: string;
  role: string;
};

export type Credentials = {
  userId: string;
  role: string;
};
