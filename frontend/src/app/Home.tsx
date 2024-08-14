import { useEffect, useRef, useState } from "react";
import { FcGoogle } from "react-icons/fc";
import { useLocation, useNavigate } from "react-router-dom";
import {
  checkSession,
  getGoogleLoginUrl,
  signout,
  verifyGoogleLogin,
} from "../api/auth";
import { DIRECT } from "../constant/constant";

function Home() {
  const googleLoginUrl = useRef("");
  const [serviceTicket, setServiceTicket] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);
  const serviceUrl = queryParams.get("service");
  const state = queryParams.get("state");
  const code = queryParams.get("code");

  useEffect(() => {
    if (code && state) {
      (async () => {
        try {
          const res = await verifyGoogleLogin(code, state);
          setLoading(false);

          if (res instanceof Error) {
            return setError(res.message);
          }

          navigate("/");
          if (state !== DIRECT)
            window.location.href = `${state}?ticket=${res.serviceTicket}`;
        } catch {
          return setError("Failed to verify Google login");
        }
      })();
    } else if (!state || !code) {
      (async () => {
        try {
          const res = await getGoogleLoginUrl(serviceUrl);
          setLoading(false);

          if (res instanceof Error) {
            return setError(res.message);
          }

          googleLoginUrl.current = res;
        } catch {
          return setError("Failed to get Google login URL");
        }
      })();

      (async () => {
        try {
          const res = await checkSession(serviceUrl);
          setLoading(false);

          if (res instanceof Error) {
            return;
          }

          setServiceTicket(res.serviceTicket);
          if (serviceUrl && res.serviceTicket)
            window.location.href = `${serviceUrl}?ticket=${res.serviceTicket}`;
        } catch {
          return setError("Failed to check session");
        }
      })();
    }
  }, [serviceUrl, state, code, navigate]);

  const handleClick = () => {
    if (loading) return;
    window.location.href = googleLoginUrl.current;
  };

  const handleSignout = async () => {
    if (!serviceTicket) return;

    try {
      await signout();
      setServiceTicket("");
    } catch (error) {
      console.error("Failed to logout: ", error);
    }
  };

  const SSOLoginStatus = () => {
    if (serviceTicket)
      return (
        <>
          {/* <h3 className="mt-4 text-2xl font-medium">Logged in as</h3>
        <p className="mt-1">{credentials.email}</p> */}
          <h3 className="mt-4 text-2xl font-medium">Logged in</h3>
          <button
            onClick={handleSignout}
            className="mt-8 flex w-[80%] items-center justify-center rounded-lg border border-gray-300 py-2 text-lg text-gray-600 duration-300 ease-in-out hover:bg-slate-100"
          >
            <h3>Logout</h3>
          </button>
        </>
      );

    return (
      <button
        onClick={handleClick}
        className="mt-8 flex items-center justify-center rounded-lg border border-gray-300 px-[30%] py-2 text-lg text-gray-600 duration-300 ease-in-out hover:bg-slate-100"
      >
        <FcGoogle className="mr-2 inline h-8 w-8" /> <h3>Google</h3>
      </button>
    );
  };

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-gray-300">
      <div className="flex h-[40vh] w-[60vw] flex-col items-center rounded-xl bg-white px-8 py-[10vh] drop-shadow-xl md:w-[40vw] xl:w-[30vw] 2xl:w-[20vw]">
        <h1 className="text-4xl font-bold">SSO Login</h1>
        {SSOLoginStatus()}
        {loading && <p className="mt-4 text-gray-500">Loading...</p>}
        {error && <p className="mt-4 text-red-500">{error}</p>}
      </div>
    </div>
  );
}

export default Home;
