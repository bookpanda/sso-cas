export type AuthenticateSSODTO = {
  access_token: string;
  refresh_token: string;
  expires_in: Date;
};

export type Credentials = {
  user_id: string;
  role: string;
};
