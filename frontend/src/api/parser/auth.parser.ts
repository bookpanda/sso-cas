import {
  CheckSessionDTO,
  Credentials,
  ServiceTicket,
  ValidateDTO,
  VerifyGoogleLoginDTO,
} from "../dto/auth.dto";

export const parseCheckSession = (dto: CheckSessionDTO): ServiceTicket => {
  return {
    serviceTicket: dto.service_ticket,
  };
};

export const parseVerifyGoogleLogin = (
  dto: VerifyGoogleLoginDTO
): ServiceTicket => {
  return {
    serviceTicket: dto.service_ticket,
  };
};

export const parseCredentials = (dto: ValidateDTO): Credentials => {
  return {
    userId: dto.user_id,
    email: dto.email,
    role: dto.role,
  };
};
