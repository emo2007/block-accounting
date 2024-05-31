import axios from "axios";
import Cookies from "js-cookie";
// const tokenStorage = Cookies.get("accessToken");

export const globalService = axios.create({
  baseURL: "http://209.141.36.222:8085/",
  timeout: 10000,
});

globalService.interceptors.request.use(
  (config) => {
    // Get the access token from the cookie
    const token = Cookies.get("accessToken");
    console.log(token);

    // If the token exists, set it in the Authorization header
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    // If there's an error in the request configuration, you can handle it here
    return Promise.reject(error);
  }
);
