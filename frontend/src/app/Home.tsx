import { FcGoogle } from "react-icons/fc";
import { useLocation } from "react-router-dom";
import { signout } from "../api/auth";
import { useAuthSSO } from "../hooks/useAuthSSO";
import { useGetGoogleLogin } from "../hooks/useGetGoogleLogin";

function Home() {
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const serviceUrl = queryParams.get("service");
  const state = queryParams.get("state");
  const code = queryParams.get("code");

  const {
    googleLoginUrl,
    loading: ggLoading,
    error: ggError,
  } = useGetGoogleLogin(serviceUrl);

  const { setServiceTicket, serviceTicket, credentials, loading, error } =
    useAuthSSO(code, state, serviceUrl);

  const handleClick = () => {
    if (ggLoading) return;
    window.location.href = googleLoginUrl;
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
          <h3 className="mt-4 text-2xl font-medium">Logged in as</h3>
          <p className="mt-1">{credentials.email}</p>
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
        {error && <p className="mt-4 text-red-500">{error.message}</p>}
        {ggError && <p className="mt-4 text-red-500">{ggError.message}</p>}
      </div>
    </div>
  );
}

export default Home;
