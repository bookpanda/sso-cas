import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { authenticateSSO } from "./api/auth";
import { SERVICE, SSO_URL, WEB_URL } from "./constant/constant";

function Home() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const serviceTicket = queryParams.get("ticket");

  useEffect(() => {
    if (serviceTicket) {
      (async () => {
        try {
          const res = await authenticateSSO(serviceTicket);
          setLoading(false);

          if (res instanceof Error) {
            return setError(res.message);
          }

          console.log(res);
        } catch {
          return setError("Failed to verify Google login");
        }
      })();
    }
  }, [serviceTicket]);

  const handleClick = () => {
    window.location.href = `${SSO_URL}?service=${WEB_URL}`;
  };

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-gray-50">
      <div className="flex h-[40vh] w-[60vw] flex-col items-center rounded-xl bg-white px-8 py-[10vh] drop-shadow-xl md:w-[40vw] xl:w-[30vw] 2xl:w-[20vw]">
        <h1 className="text-4xl font-bold">{SERVICE}</h1>
        <button
          onClick={handleClick}
          className="mt-8 flex items-center justify-center rounded-lg border border-gray-300 px-[30%] py-2 text-lg text-gray-600 duration-300 ease-in-out hover:bg-slate-100"
        >
          <h3>Login via SSO</h3>
        </button>
        {loading && <p className="mt-4 text-gray-500">Loading...</p>}
        {error && <p className="mt-4 text-red-500">{error}</p>}
      </div>
    </div>
  );
}

export default Home;
