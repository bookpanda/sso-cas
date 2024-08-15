import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { checkSession, validateST, verifyGoogleLogin } from "../api/auth";
import { Credentials } from "../api/dto/auth.dto";
import { DIRECT } from "../constant/constant";

interface AuthSSOPayload {
  setServiceTicket: React.Dispatch<React.SetStateAction<string>>;
  serviceTicket: string;
  credentials: Credentials;
  loading: boolean;
  error: Error | null;
}

export const useAuthSSO = (
  code: string | null,
  state: string | null,
  serviceUrl: string | null
): AuthSSOPayload => {
  const [serviceTicket, setServiceTicket] = useState("");
  const [credentials, setCredentials] = useState<Credentials>({
    userId: "",
    email: "",
    role: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    setLoading(true);
    if (code && state)
      (async () => {
        const res = await verifyGoogleLogin(code, state);
        if (res instanceof Error) {
          return setError(res);
        }
        setServiceTicket(res.serviceTicket);

        if (state !== DIRECT)
          // is login from other service
          window.location.href = `${state}?ticket=${res.serviceTicket}`;
        else navigate("/");
      })();
    else
      (async () => {
        const res = await checkSession(serviceUrl);
        if (res instanceof Error) {
          return setError(res);
        }
        setServiceTicket(res.serviceTicket);

        if (serviceUrl && res.serviceTicket)
          // service came to validate session, redirecting back
          window.location.href = `${serviceUrl}?ticket=${res.serviceTicket}`;
        else {
          const res2 = await validateST(res.serviceTicket);
          if (res2 instanceof Error) {
            return;
          }
          setCredentials(res2);
        }
      })();
    setLoading(false);
  }, [code, state, serviceUrl, navigate]);

  return { setServiceTicket, serviceTicket, credentials, loading, error };
};
