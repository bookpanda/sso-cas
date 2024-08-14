import {
  CheckSessionDTO,
  ServiceTicket,
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
