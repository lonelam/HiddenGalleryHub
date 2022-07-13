import { Axios } from "axios";

export const apiClient = new Axios({
  responseType: "json",
  transformResponse: (data) => {
    try {
      return JSON.parse(data);
    } catch {
      return data;
    }
  },
});

export default apiClient;
