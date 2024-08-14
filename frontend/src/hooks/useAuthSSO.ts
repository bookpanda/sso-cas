import { useEffect, useState } from "react";
import { checkSession, verifyGoogleLogin } from "../api/auth";

interface AuthSSOPayload {
  serviceTicket: string;
  loading: boolean;
  error: Error | null;
}

export const useAuthSSO = (
  code: string | null,
  state: string | null,
  serviceUrl: string | null
): AuthSSOPayload => {
  const [serviceTicket, setServiceTicket] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    setLoading(true);
    if (code && state)
      (async () => {
        const res = await verifyGoogleLogin(code, state);

        if (res instanceof Error) {
          return setError(res);
        }

        setServiceTicket(res.serviceTicket);
      })();
    else
      (async () => {
        const res = await checkSession(serviceUrl);

        if (res instanceof Error) {
          return setError(res);
        }

        setServiceTicket(res.serviceTicket);
      })();
    setLoading(false);
  }, [code, state, serviceUrl]);

  return { serviceTicket, loading, error };
};
