import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { authenticateSSO, logout } from "./api/auth";
import { SERVICE, SSO_URL, WEB_URL } from "./constant/constant";

function Home() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);
  const serviceTicket = queryParams.get("ticket");

  const [accessToken, setAccessToken] = useState(
    localStorage.getItem("access_token") || ""
  );

  useEffect(() => {
    if (serviceTicket) {
      (async () => {
        try {
          const res = await authenticateSSO(serviceTicket);
          setLoading(false);

          if (res instanceof Error) {
            return setError(res.message);
          }

          localStorage.setItem("access_token", res.accessToken);
          localStorage.setItem("refresh_token", res.refreshToken);
          setAccessToken(res.accessToken);
          navigate("/");
          console.log(res);
        } catch {
          return setError("Failed to verify Google login");
        }
      })();
    }
  }, [serviceTicket, navigate]);

  const handleClick = () => {
    window.location.href = `${SSO_URL}?service=${WEB_URL}`;
  };

  const handleLogout = async () => {
    if (!accessToken) return;

    try {
      await logout(accessToken);
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      setAccessToken("");
    } catch (error) {
      console.error("Failed to logout: ", error);
    }
  };

  const SSOLoginStatus = () => {
    if (accessToken)
      return (
        <>
          <h3 className="mt-4 text-2xl font-medium">Logged in</h3>
          <button
            onClick={handleLogout}
            className="mt-8 flex w-[80%] items-center justify-center rounded-lg border border-gray-300 py-2 text-lg text-gray-600 duration-300 ease-in-out hover:bg-slate-100"
          >
            <h3>Logout</h3>
          </button>
        </>
      );

    return (
      <button
        onClick={handleClick}
        className="mt-8 flex w-[80%] items-center justify-center rounded-lg border border-gray-300 py-2 text-lg text-gray-600 duration-300 ease-in-out hover:bg-slate-100"
      >
        <h3>Login via SSO</h3>
      </button>
    );
  };

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-gray-50">
      <div className="flex h-[40vh] w-[60vw] flex-col items-center rounded-xl bg-white px-8 py-[10vh] drop-shadow-xl md:w-[40vw] xl:w-[30vw] 2xl:w-[20vw]">
        <h1 className="text-4xl font-bold">{SERVICE}</h1>
        {SSOLoginStatus()}

        {loading && <p className="mt-4 text-gray-500">Loading...</p>}
        {error && <p className="mt-4 text-red-500">{error}</p>}
      </div>
    </div>
  );
}

export default Home;
