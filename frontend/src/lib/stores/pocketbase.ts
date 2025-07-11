import PocketBase from "pocketbase";
import { writable } from "svelte/store";
import { goto } from "$app/navigation";
import type { Permission } from "$lib/types/permission";

// Use environment variable for backend URL, fallback to localhost for development
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8090';

// connect to backend
const pb = new PocketBase(BACKEND_URL);
pb.autoCancellation(false);

// Add response interceptor for auth-related errors
pb.beforeSend = function (url: string, options: Record<string, any>) {
  const originalResponse = options.fetch || fetch;
  options.fetch = async (input: RequestInfo, init?: RequestInit) => {
    const response = await originalResponse(input, init);

    if (!response.ok) {
      // Handle different error cases
      switch (response.status) {
        case 401:
          pb.authStore.clear();
          goto("/authentication/sign-in");
          break;
        case 403:
          try {
            const data = await response.clone().json();
            if (data.message === "Password change required") {
              goto("/change-password");
            }
          } catch (e) {
            console.error("Error parsing 403 response:", e);
          }
          break;
      }
    }

    return response;
  };
  return { url, options };
};

// Try to load auth state from storage if we're in a browser environment
if (typeof window !== "undefined") {
  try {
    const storedToken = localStorage.getItem("pocketbase_auth");
    if (storedToken) {
      const authData = JSON.parse(storedToken);
      pb.authStore.save(authData.token, authData.model);
    }
  } catch (error) {
    console.error("Error loading stored auth:", error);
  }
}

// export stores
export const pocketbase = writable(pb);
export const permission = writable({} as Permission);

// Export the instance for direct usage
export default pb;

// Export the backend URL for use in other parts of the app
export const getBackendUrl = () => BACKEND_URL;
