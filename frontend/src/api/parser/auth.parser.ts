import { VerifyGoogleLoginDTO } from "../dto/auth.dto";

export const parseVerifyGoogleLogin = (dto: VerifyGoogleLoginDTO) => {
  return {
    serviceTicket: dto.service_ticket,
  };
};
