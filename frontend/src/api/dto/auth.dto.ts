export type CheckSessionDTO = {
  service_ticket: string;
};

export type VerifyGoogleLoginDTO = {
  service_ticket: string;
};

export type ServiceTicket = {
  serviceTicket: string;
};

export type ValidateDTO = {
  user_id: string;
  email: string;
  role: string;
};

export type Credentials = {
  userId: string;
  email: string;
  role: string;
};
