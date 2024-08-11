import { useEffect, useRef, useState } from "react";
import { FcGoogle } from "react-icons/fc";
import { getGoogleLoginUrl } from "./api/auth";

function App() {
  const googleLoginUrl = useRef("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchGoogleLogin = async () => {
      try {
        const res = await getGoogleLoginUrl("http://localhost:3000");
        setLoading(false);

        if (res instanceof Error) {
          return setError(res.message);
        }

        googleLoginUrl.current = res.data;
      } catch {
        return setError("Failed to get Google login URL");
      }
    };

    fetchGoogleLogin();
  }, []);

  const handleClick = () => {
    if (loading) return;
    window.location.href = googleLoginUrl.current;
  };

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-gray-50">
      <div className="flex h-[40vh] w-[60vw] flex-col items-center rounded-xl bg-white px-8 py-[10vh] drop-shadow-xl md:w-[40vw] xl:w-[30vw] 2xl:w-[20vw]">
        <h1 className="text-4xl font-bold">SSO Login</h1>
        <button
          onClick={handleClick}
          className="mt-8 flex items-center justify-center rounded-lg border border-gray-300 px-[30%] py-2 text-lg text-gray-600 duration-300 ease-in-out hover:bg-slate-100"
        >
          <FcGoogle className="mr-2 inline h-8 w-8" /> <h3>Google</h3>
        </button>
        {loading && <p className="mt-4 text-gray-500">Loading...</p>}
        {error && <p className="mt-4 text-red-500">{error}</p>}
      </div>
    </div>
  );
}

export default App;
