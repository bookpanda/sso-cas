import { useEffect, useState } from "react";
import { getGoogleLoginUrl } from "../api/auth";

interface GetGoogleLoginPayload {
  googleLoginUrl: string;
  loading: boolean;
  error: Error | null;
}

export const useGetGoogleLogin = (
  serviceUrl: string | null
): GetGoogleLoginPayload => {
  const [googleLoginUrl, setGoogleLoginUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    setLoading(true);
    (async () => {
      const res = await getGoogleLoginUrl(serviceUrl);

      if (res instanceof Error) {
        return setError(res);
      }

      setGoogleLoginUrl(res);
    })();
    setLoading(false);
  }, [serviceUrl]);

  return { googleLoginUrl: googleLoginUrl, loading, error };
};
