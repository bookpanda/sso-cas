import { CheckSessionDTO, VerifyGoogleLoginDTO } from "../dto/auth.dto";

export const parseCheckSession = (dto: CheckSessionDTO) => {
  return {
    serviceTicket: dto.service_ticket,
  };
};

export const parseVerifyGoogleLogin = (dto: VerifyGoogleLoginDTO) => {
  return {
    serviceTicket: dto.service_ticket,
  };
};
