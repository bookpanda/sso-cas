import { FcGoogle } from "react-icons/fc";
import { useLocation } from "react-router-dom";
import { signout } from "../api/auth";
import { Button } from "../components/Button";
import { SERVICES_URL } from "../constant/constant";
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
          <Button onClick={handleSignout}>Logout</Button>
          <div className="flex w-[80%] space-x-4">
            {SERVICES_URL.map((url, idx) => (
              <Button key={url} onClick={() => window.open(url, "_blank")}>
                Service {idx + 1}
              </Button>
            ))}
          </div>
        </>
      );

    return (
      <Button onClick={handleClick}>
        <FcGoogle className="mr-2 inline h-8 w-8" /> <h3>Google</h3>
      </Button>
    );
  };

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-gray-300">
      <div className="flex h-[45vh] w-[60vw] flex-col items-center justify-center rounded-xl bg-white px-8 py-[10vh] drop-shadow-xl md:w-[40vw] xl:w-[30vw] 2xl:w-[20vw]">
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
